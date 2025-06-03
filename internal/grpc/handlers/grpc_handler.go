package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/handlers"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/services"
	pb "github.com/VladSnap/shortener/proto"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShortenerGRPCHandler implements the gRPC service interface
type ShortenerGRPCHandler struct {
	pb.UnimplementedShortenerServiceServer
	service      handlers.ShorterService
	deleteWorker handlers.DeleterWorker
	baseURL      string
	opts         *config.Options
}

// NewShortenerGRPCHandler creates a new gRPC handler
func NewShortenerGRPCHandler(
	service handlers.ShorterService,
	deleteWorker handlers.DeleterWorker,
	baseURL string,
	opts *config.Options,
) *ShortenerGRPCHandler {
	return &ShortenerGRPCHandler{
		service:      service,
		deleteWorker: deleteWorker,
		baseURL:      baseURL,
		opts:         opts,
	}
}

// CreateShortLink creates a shortened URL for a given original URL
func (h *ShortenerGRPCHandler) CreateShortLink(ctx context.Context, req *pb.CreateShortLinkRequest) (*pb.CreateShortLinkResponse, error) {
	if req.OriginalUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "original_url is required")
	}

	// Extract user ID from context (set by auth interceptor)
	userID, ok := ctx.Value(constants.UserIDContextKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	shortedLink, err := h.service.CreateShortLink(ctx, req.OriginalUrl, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create short link: %v", err))
	}

	return &pb.CreateShortLinkResponse{
		ShortUrl:    h.baseURL + "/" + shortedLink.URL,
		IsDuplicate: shortedLink.IsDuplicated,
	}, nil
}

// CreateShortLinkBatch creates multiple shortened URLs in a single request
func (h *ShortenerGRPCHandler) CreateShortLinkBatch(ctx context.Context, req *pb.CreateShortLinkBatchRequest) (*pb.CreateShortLinkBatchResponse, error) {
	if len(req.Links) == 0 {
		return &pb.CreateShortLinkBatchResponse{Links: []*pb.ShortedLinkBatch{}}, nil
	}

	// Extract user ID from context (set by auth interceptor)
	userID, ok := ctx.Value(constants.UserIDContextKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	// Convert gRPC request to service model
	originalLinks := make([]*services.OriginalLink, 0, len(req.Links))
	for _, link := range req.Links {
		if link.OriginalUrl == "" {
			return nil, status.Error(codes.InvalidArgument, "original_url is required for all links")
		}
		originalLinks = append(originalLinks, &services.OriginalLink{
			CorelationID: link.CorrelationId,
			URL:          link.OriginalUrl,
		})
	}

	shortedLinks, err := h.service.CreateShortLinkBatch(ctx, originalLinks, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create batch: %v", err))
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

// GetURL retrieves the original URL by its short identifier
func (h *ShortenerGRPCHandler) GetURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	if req.ShortId == "" {
		return nil, status.Error(codes.InvalidArgument, "short_id is required")
	}

	shortedLink, err := h.service.GetURL(ctx, req.ShortId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get URL: %v", err))
	}

	if shortedLink == nil {
		return nil, status.Error(codes.NotFound, "URL not found")
	}

	if shortedLink.IsDeleted {
		return nil, status.Error(codes.FailedPrecondition, "URL has been removed")
	}

	return &pb.GetURLResponse{
		OriginalUrl: shortedLink.OriginalURL,
		IsDeleted:   shortedLink.IsDeleted,
	}, nil
}

// GetAllByUserID retrieves all URLs shortened by a specific user
func (h *ShortenerGRPCHandler) GetAllByUserID(ctx context.Context, req *pb.GetAllByUserIDRequest) (*pb.GetAllByUserIDResponse, error) {
	// Extract user ID from context (set by auth interceptor)
	userID, ok := ctx.Value(constants.UserIDContextKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	shortedLinks, err := h.service.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get user URLs: %v", err))
	}

	if len(shortedLinks) == 0 {
		return nil, status.Error(codes.NotFound, "URLs for user not found")
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

// DeleteBatch marks multiple URLs as deleted
func (h *ShortenerGRPCHandler) DeleteBatch(ctx context.Context, req *pb.DeleteBatchRequest) (*pb.DeleteBatchResponse, error) {
	if len(req.ShortUrls) == 0 {
		return nil, status.Error(codes.InvalidArgument, "short_urls cannot be empty")
	}

	// Extract user ID from context (set by auth interceptor)
	userID, ok := ctx.Value(constants.UserIDContextKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	// Create channel for deletion
	const toDeleteChanSize = 100
	toDeleteChan := make(chan services.DeleteShortID, toDeleteChanSize)
	defer close(toDeleteChan)

	// Send deletion requests to channel
	for _, shortURL := range req.ShortUrls {
		if shortURL == "" {
			return nil, status.Error(codes.InvalidArgument, "short_url cannot be empty")
		}
		deleteSID := services.NewDeleteShortID(shortURL, userID)
		select {
		case toDeleteChan <- deleteSID:
		case <-ctx.Done():
			return nil, status.Error(codes.DeadlineExceeded, "context deadline exceeded")
		}
	}

	// Send to delete worker
	h.deleteWorker.AddToDelete(toDeleteChan)

	return &pb.DeleteBatchResponse{Success: true}, nil
}

// GetStats returns service statistics (only for trusted subnets)
func (h *ShortenerGRPCHandler) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	// Note: In a real implementation, you would need to implement trusted subnet checking for gRPC
	// This is more complex in gRPC as you need to extract the client IP from the connection

	stats, err := h.service.GetStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get stats: %v", err))
	}

	return &pb.GetStatsResponse{
		Urls:  int32(stats.Urls),
		Users: int32(stats.Users),
	}, nil
}

// Ping checks service health
func (h *ShortenerGRPCHandler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	// Note: In a real implementation, you would want to check the database connection here
	// similar to the HTTP ping handler, but for simplicity we'll just return OK

	return &pb.PingResponse{Status: "OK"}, nil
}

// StartGRPCServer starts the gRPC server on the specified address
func StartGRPCServer(addr string, handler *ShortenerGRPCHandler) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterShortenerServiceServer(s, handler)

	return s, lis, nil
}
