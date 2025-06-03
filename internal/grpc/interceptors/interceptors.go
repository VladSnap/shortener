// Package interceptors provides gRPC interceptors for authentication, logging, and trusted subnet validation.
package interceptors

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/VladSnap/shortener/internal/auth"
	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// AuthInterceptor provides authentication functionality for gRPC.
func AuthInterceptor(opts *config.Options) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			// No metadata, create new user ID
			userID := generateNewUserID()
			ctx = context.WithValue(ctx, constants.UserIDContextKey, userID)
			return handler(ctx, req)
		}

		// Check for auth cookie in metadata
		authCookies := md.Get("auth-cookie")
		if len(authCookies) == 0 {
			// No auth cookie, create new user ID
			userID := generateNewUserID()
			ctx = context.WithValue(ctx, constants.UserIDContextKey, userID)
			return handler(ctx, req)
		}

		authCookie := authCookies[0]

		// Verify the signed cookie
		if _, err := auth.VerifySignCookie(authCookie, opts.AuthCookieKey); err != nil {
			log.Zap.Warn("failed to verify gRPC auth cookie", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "invalid authentication")
		}

		// Decode the cookie to get user ID
		authData, err := auth.DecodeCookie(authCookie)
		if err != nil {
			log.Zap.Warn("failed to decode gRPC auth cookie", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "invalid authentication data")
		}

		// Add user ID to context
		ctx = context.WithValue(ctx, constants.UserIDContextKey, authData.UserID)
		return handler(ctx, req)
	}
}

// LoggingInterceptor provides logging functionality for gRPC.
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Extract metadata for logging
		md, _ := metadata.FromIncomingContext(ctx)
		metadataStr := formatMetadata(md)

		// Extract peer info
		peerInfo, _ := peer.FromContext(ctx)
		clientAddr := "unknown"
		if peerInfo != nil {
			clientAddr = peerInfo.Addr.String()
		}

		// Call the handler
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		// Determine status
		grpcStatus := codes.OK
		if err != nil {
			if s, ok := status.FromError(err); ok {
				grpcStatus = s.Code()
			} else {
				grpcStatus = codes.Internal
			}
		}

		// Log the request
		log.Zap.Info("gRPC Request",
			zap.String("method", info.FullMethod),
			zap.String("client_addr", clientAddr),
			zap.String("status", grpcStatus.String()),
			zap.Duration("duration", duration),
			zap.String("metadata", metadataStr),
			zap.Error(err),
		)

		return resp, err
	}
}

// TrustedSubnetInterceptor provides trusted subnet validation for gRPC.
func TrustedSubnetInterceptor(trustedSubnet string) grpc.UnaryServerInterceptor {
	_, subnet, _ := net.ParseCIDR(trustedSubnet)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Only apply to GetStats method
		if !strings.HasSuffix(info.FullMethod, "GetStats") {
			return handler(ctx, req)
		}

		if subnet == nil {
			return nil, status.Error(codes.PermissionDenied, "trusted subnet not configured")
		}

		// Extract client IP from peer info
		peerInfo, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "unable to get client address")
		}

		// Parse the client address
		host, _, err := net.SplitHostPort(peerInfo.Addr.String())
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "invalid client address")
		}

		// Check if we have X-Real-IP in metadata (for proxy scenarios)
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			realIPs := md.Get("x-real-ip")
			if len(realIPs) > 0 {
				host = realIPs[0]
			}
		}

		ip := net.ParseIP(host)
		if ip == nil {
			return nil, status.Error(codes.PermissionDenied, "invalid IP address")
		}

		if !subnet.Contains(ip) {
			return nil, status.Error(codes.PermissionDenied, "IP address not in trusted subnet")
		}

		return handler(ctx, req)
	}
}

// generateNewUserID creates a new UUID for unauthenticated users.
func generateNewUserID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Zap.Error("failed to generate new user ID", zap.Error(err))
		return "anonymous"
	}
	return id.String()
}

// formatMetadata converts gRPC metadata to a readable string.
func formatMetadata(md metadata.MD) string {
	if md == nil {
		return ""
	}

	parts := make([]string, 0, len(md))
	for k, v := range md {
		parts = append(parts, fmt.Sprintf("%s: %v", k, v))
	}
	return strings.Join(parts, " | ")
}
