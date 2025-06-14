package handlers

import (
	"context"
	"fmt"
	"net"

	"github.com/VladSnap/shortener/internal/config"
	grpcvalidation "github.com/VladSnap/shortener/internal/grpc/validation"
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

	// Error message formats for wrapping.
	validationErrorFormat     = "validation failed: %w"
	userExtractionErrorFormat = "user extraction failed: %w"
	contextErrorFormat        = "context validation failed: %w"
	urlLookupErrorFormat      = "URL lookup failed: %w"
	urlAccessErrorFormat      = "URL access failed: %w"
	userURLsLookupErrorFormat = "user URLs lookup failed: %w"
	contextDeadlineFormat     = "context deadline exceeded: %w"
)

// handleServiceError обрабатывает ошибки сервиса и возвращает соответствующую gRPC ошибку.
func handleServiceError(err error, operation string) error {
	if err == nil {
		return nil
	}
	return status.Errorf(codes.Internal, "failed to %s: %v", operation, err)
}

// handleDatabaseError обрабатывает ошибки базы данных.
func handleDatabaseError(err error, message string) error {
	if err == nil {
		return nil
	}
	return status.Errorf(codes.Unavailable, "%s: %v", message, err)
}

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
	if err := grpcvalidation.ValidateOriginalURL(req.GetOriginalUrl()); err != nil {
		return nil, fmt.Errorf(validationErrorFormat, err)
	}

	userID, err := grpcvalidation.ExtractUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf(userExtractionErrorFormat, err)
	}

	shortedLink, err := h.service.CreateShortLink(ctx, req.GetOriginalUrl(), userID)
	if err != nil {
		return nil, handleServiceError(err, "create short link")
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

	userID, err := grpcvalidation.ExtractUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf(userExtractionErrorFormat, err)
	}

	// Convert gRPC request to service model
	originalLinks := make([]*services.OriginalLink, 0, len(req.GetLinks()))
	for _, link := range req.GetLinks() {
		if err := grpcvalidation.ValidateOriginalURL(link.GetOriginalUrl()); err != nil {
			return nil, fmt.Errorf(validationErrorFormat, err)
		}

		originalLinks = append(originalLinks, &services.OriginalLink{
			CorelationID: link.GetCorrelationId(),
			URL:          link.GetOriginalUrl(),
		})
	}

	shortedLinks, err := h.service.CreateShortLinkBatch(ctx, originalLinks, userID)
	if err != nil {
		return nil, handleServiceError(err, "create batch")
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
	if err := grpcvalidation.ValidateShortID(req.GetShortId()); err != nil {
		return nil, fmt.Errorf(validationErrorFormat, err)
	}

	shortedLink, err := h.service.GetURL(ctx, req.GetShortId())
	if err != nil {
		return nil, handleServiceError(err, "get URL")
	}

	if shortedLink == nil {
		return nil, fmt.Errorf(urlLookupErrorFormat, status.Error(codes.NotFound, "URL not found"))
	}

	if shortedLink.IsDeleted {
		return nil, fmt.Errorf(urlAccessErrorFormat, status.Error(codes.FailedPrecondition, "URL has been removed"))
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
	userID, err := grpcvalidation.ExtractUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf(userExtractionErrorFormat, err)
	}

	shortedLinks, err := h.service.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, handleServiceError(err, "get user URLs")
	}

	if len(shortedLinks) == 0 {
		return nil, fmt.Errorf(userURLsLookupErrorFormat, status.Error(codes.NotFound, "URLs for user not found"))
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
	if err := grpcvalidation.ValidateShortURLs(req.GetShortUrls()); err != nil {
		return nil, fmt.Errorf(validationErrorFormat, err)
	}

	userID, err := grpcvalidation.ExtractUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf(userExtractionErrorFormat, err)
	}

	// Create channel for deletion
	toDeleteChan := make(chan services.DeleteShortID, toDeleteChanSize)
	defer close(toDeleteChan)

	// Send deletion requests to channel
	for _, shortURL := range req.GetShortUrls() {
		if err := grpcvalidation.ValidateContextDeadline(ctx); err != nil {
			return nil, fmt.Errorf(contextErrorFormat, err)
		}

		deleteSID := services.NewDeleteShortID(shortURL, userID)
		select {
		case toDeleteChan <- deleteSID:
		case <-ctx.Done():
			return nil, fmt.Errorf(contextDeadlineFormat,
				status.Error(codes.DeadlineExceeded, "context deadline exceeded"))
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
		return nil, handleServiceError(err, "get stats")
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
		return nil, handleDatabaseError(err, "database not available")
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
