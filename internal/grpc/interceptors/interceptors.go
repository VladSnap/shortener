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

// TrustedSubnetConfig holds configuration for trusted subnet validation.
type TrustedSubnetConfig struct {
	// TrustedSubnet is the CIDR notation of the trusted subnet
	TrustedSubnet string
	// ProtectedMethods is a list of gRPC method names that require trusted subnet validation
	// If empty, validation is applied to all methods
	ProtectedMethods []string
	// UseMethodSuffix when true, matches method names by suffix instead of exact match
	UseMethodSuffix bool
}

// TrustedSubnetInterceptor provides configurable trusted subnet validation for gRPC.
func TrustedSubnetInterceptor(config TrustedSubnetConfig) grpc.UnaryServerInterceptor {
	_, subnet, err := net.ParseCIDR(config.TrustedSubnet)
	if err != nil {
		log.Zap.Error("failed to parse trusted subnet CIDR",
			zap.String("subnet", config.TrustedSubnet),
			zap.Error(err))
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Check if this method requires trusted subnet validation
		if !shouldValidateMethod(info.FullMethod, config) {
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
			log.Zap.Warn("access denied: IP not in trusted subnet",
				zap.String("method", info.FullMethod),
				zap.String("client_ip", ip.String()),
				zap.String("trusted_subnet", config.TrustedSubnet))
			return nil, status.Error(codes.PermissionDenied, "IP address not in trusted subnet")
		}

		return handler(ctx, req)
	}
}

// shouldValidateMethod determines if a method requires trusted subnet validation.
func shouldValidateMethod(fullMethod string, config TrustedSubnetConfig) bool {
	// If no protected methods are specified, validate all methods
	if len(config.ProtectedMethods) == 0 {
		return true
	}

	// Check each protected method
	for _, protectedMethod := range config.ProtectedMethods {
		if config.UseMethodSuffix {
			if strings.HasSuffix(fullMethod, protectedMethod) {
				return true
			}
		} else {
			if fullMethod == protectedMethod || strings.HasSuffix(fullMethod, "/"+protectedMethod) {
				return true
			}
		}
	}

	return false
}

// NewTrustedSubnetConfigWithSuffix creates a new configuration that matches methods by suffix.
func NewTrustedSubnetConfigWithSuffix(trustedSubnet string, methodSuffixes ...string) TrustedSubnetConfig {
	return TrustedSubnetConfig{
		TrustedSubnet:    trustedSubnet,
		ProtectedMethods: methodSuffixes,
		UseMethodSuffix:  true,
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
