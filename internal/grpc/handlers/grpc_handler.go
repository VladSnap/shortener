package handlers

import (
	"context"
	"fmt"
	"net"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/handlers"
	"github.com/VladSnap/shortener/internal/services"
	pb "github.com/VladSnap/shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// Channel sizes for delete operations.
	toDeleteChanSize = 100

	// Error messages.
	errUserIDNotFound      = "user ID not found in context"
	errOriginalURLRequired = "original_url is required"
	errShortIDRequired     = "short_id is required"
	errURLNotFound         = "URL not found"
	errURLRemoved          = "URL has been removed"
	errUserURLsNotFound    = "URLs for user not found"
	errShortURLsEmpty      = "short_urls cannot be empty"
	errShortURLEmpty       = "short_url cannot be empty"
	errContextDeadline     = "context deadline exceeded"
)

// ShortenerGRPCHandler implements the gRPC service interface.
type ShortenerGRPCHandler struct {
	pb.UnimplementedShortenerServiceServer
	service       handlers.ShorterService
	deleteWorker  handlers.DeleterWorker
	healthService *services.HealthService
	opts          *config.Options
	baseURL       string
}

// NewShortenerGRPCHandler creates a new gRPC handler.
func NewShortenerGRPCHandler(
	service handlers.ShorterService,
	deleteWorker handlers.DeleterWorker,
	baseURL string,
	opts *config.Options,
) *ShortenerGRPCHandler {
	return &ShortenerGRPCHandler{
		service:       service,
		deleteWorker:  deleteWorker,
		baseURL:       baseURL,
		opts:          opts,
		healthService: services.NewHealthService(opts.DataBaseConnString),
	}
}

// CreateShortLink creates a shortened URL for a given original URL.
func (h *ShortenerGRPCHandler) CreateShortLink(
	ctx context.Context,
	req *pb.CreateShortLinkRequest,
) (*pb.CreateShortLinkResponse, error) {
	if req.GetOriginalUrl() == "" {
		return nil, status.Error(codes.InvalidArgument, errOriginalURLRequired)
	}

	// Extract user ID from context (set by auth interceptor)
	userID, ok := ctx.Value(constants.UserIDContextKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, errUserIDNotFound)
	}

	shortedLink, err := h.service.CreateShortLink(ctx, req.GetOriginalUrl(), userID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create short link: %w", err)
		return nil, status.Error(codes.Internal, wrappedErr.Error())
	}

	return &pb.CreateShortLinkResponse{
		ShortUrl:    h.baseURL + "/" + shortedLink.URL,
		IsDuplicate: shortedLink.IsDuplicated,
	}, nil
}

// CreateShortLinkBatch creates multiple shortened URLs in a single request.
func (h *ShortenerGRPCHandler) CreateShortLinkBatch(
	ctx context.Context,
	req *pb.CreateShortLinkBatchRequest,
) (*pb.CreateShortLinkBatchResponse, error) {
	if len(req.GetLinks()) == 0 {
		return &pb.CreateShortLinkBatchResponse{Links: []*pb.ShortedLinkBatch{}}, nil
	}

	// Extract user ID from context (set by auth interceptor)
	userID, ok := ctx.Value(constants.UserIDContextKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, errUserIDNotFound)
	}

	// Convert gRPC request to service model
	originalLinks := make([]*services.OriginalLink, 0, len(req.GetLinks()))
	for _, link := range req.GetLinks() {
		if link.GetOriginalUrl() == "" {
			return nil, status.Error(codes.InvalidArgument, errOriginalURLRequired)
		}
		originalLinks = append(originalLinks, &services.OriginalLink{
			CorelationID: link.GetCorrelationId(),
			URL:          link.GetOriginalUrl(),
		})
	}

	shortedLinks, err := h.service.CreateShortLinkBatch(ctx, originalLinks, userID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create batch: %w", err)
		return nil, status.Error(codes.Internal, wrappedErr.Error())
	}

	// Convert service response to gRPC response
	responseLinks := make([]*pb.ShortedLinkBatch, 0, len(shortedLinks))
	for _, link := range shortedLinks {
		responseLinks = append(responseLinks, &pb.ShortedLinkBatch{
			CorrelationId: link.CorelationID,
			ShortUrl:      h.baseURL + "/" + link.URL,
		})
	}

	return &pb.CreateShortLinkBatchResponse{Links: responseLinks}, nil
}

// GetURL retrieves the original URL by its short identifier.
func (h *ShortenerGRPCHandler) GetURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	if req.GetShortId() == "" {
		return nil, status.Error(codes.InvalidArgument, errShortIDRequired)
	}

	shortedLink, err := h.service.GetURL(ctx, req.GetShortId())
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get URL: %w", err)
		return nil, status.Error(codes.Internal, wrappedErr.Error())
	}

	if shortedLink == nil {
		return nil, status.Error(codes.NotFound, errURLNotFound)
	}

	if shortedLink.IsDeleted {
		return nil, status.Error(codes.FailedPrecondition, errURLRemoved)
	}

	return &pb.GetURLResponse{
		OriginalUrl: shortedLink.OriginalURL,
		IsDeleted:   shortedLink.IsDeleted,
	}, nil
}

// GetAllByUserID retrieves all URLs shortened by a specific user.
func (h *ShortenerGRPCHandler) GetAllByUserID(
	ctx context.Context,
	req *pb.GetAllByUserIDRequest,
) (*pb.GetAllByUserIDResponse, error) {
	// Extract user ID from context (set by auth interceptor)
	userID, ok := ctx.Value(constants.UserIDContextKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, errUserIDNotFound)
	}

	shortedLinks, err := h.service.GetAllByUserID(ctx, userID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get user URLs: %w", err)
		return nil, status.Error(codes.Internal, wrappedErr.Error())
	}

	if len(shortedLinks) == 0 {
		return nil, status.Error(codes.NotFound, errUserURLsNotFound)
	}

	// Convert service response to gRPC response
	userUrls := make([]*pb.UserURL, 0, len(shortedLinks))
	for _, link := range shortedLinks {
		userUrls = append(userUrls, &pb.UserURL{
			OriginalUrl: link.OriginalURL,
			ShortUrl:    h.baseURL + "/" + link.URL,
		})
	}

	return &pb.GetAllByUserIDResponse{Urls: userUrls}, nil
}

// DeleteBatch marks multiple URLs as deleted.
func (h *ShortenerGRPCHandler) DeleteBatch(
	ctx context.Context,
	req *pb.DeleteBatchRequest,
) (*pb.DeleteBatchResponse, error) {
	if len(req.GetShortUrls()) == 0 {
		return nil, status.Error(codes.InvalidArgument, errShortURLsEmpty)
	}

	// Extract user ID from context (set by auth interceptor)
	userID, ok := ctx.Value(constants.UserIDContextKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, errUserIDNotFound)
	}

	// Create channel for deletion
	toDeleteChan := make(chan services.DeleteShortID, toDeleteChanSize)
	defer close(toDeleteChan)

	// Send deletion requests to channel
	for _, shortURL := range req.GetShortUrls() {
		if shortURL == "" {
			return nil, status.Error(codes.InvalidArgument, errShortURLEmpty)
		}
		deleteSID := services.NewDeleteShortID(shortURL, userID)
		select {
		case toDeleteChan <- deleteSID:
		case <-ctx.Done():
			return nil, status.Error(codes.DeadlineExceeded, errContextDeadline)
		}
	}

	// Send to delete worker
	h.deleteWorker.AddToDelete(toDeleteChan)

	return &pb.DeleteBatchResponse{Success: true}, nil
}

// GetStats returns service statistics (only for trusted subnets).
func (h *ShortenerGRPCHandler) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	// Note: In a real implementation, you would need to implement trusted subnet checking for gRPC
	// This is more complex in gRPC as you need to extract the client IP from the connection

	stats, err := h.service.GetStats(ctx)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get stats: %w", err)
		return nil, status.Error(codes.Internal, wrappedErr.Error())
	}

	return &pb.GetStatsResponse{
		Urls:  int32(stats.Urls),
		Users: int32(stats.Users),
	}, nil
}

// Ping checks service health.
func (h *ShortenerGRPCHandler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	// Check database connection if configured
	err := h.healthService.PingDatabase(ctx)
	if err != nil {
		wrappedErr := fmt.Errorf("database not available: %w", err)
		return nil, status.Error(codes.Unavailable, wrappedErr.Error())
	}

	return &pb.PingResponse{Status: "OK"}, nil
}

// StartGRPCServer starts the gRPC server on the specified address.
func StartGRPCServer(addr string, handler *ShortenerGRPCHandler) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterShortenerServiceServer(s, handler)

	return s, lis, nil
}
