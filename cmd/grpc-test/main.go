package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/VladSnap/shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// Server configuration constants.
	grpcServerAddr = "127.0.0.1:9090"
	testTimeout    = 100 * time.Second
	testURL        = "https://example.com"
)

// closeConn safely closes connection and logs any errors.
func closeConn(conn io.Closer) {
	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}

func main() {
	// Connect to gRPC server
	conn, err := grpc.NewClient(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer closeConn(conn)

	client := pb.NewShortenerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Test Ping
	log.Println("Testing Ping...")
	pingResp, err := client.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		log.Fatalf("Ping failed: %v", err)
	}
	log.Printf("Ping response: %s\n", pingResp.GetStatus())

	// Test CreateShortLink
	log.Println("Testing CreateShortLink...")

	// For testing, let the server generate a new user ID
	// In a real client, you would include a properly signed auth-cookie

	createResp, err := client.CreateShortLink(ctx, &pb.CreateShortLinkRequest{
		OriginalUrl: testURL,
	})
	if err != nil {
		log.Fatalf("CreateShortLink failed: %v", err)
	}
	log.Printf("Created short link: %s (duplicate: %t)\n", createResp.GetShortUrl(), createResp.GetIsDuplicate())

	// Test GetStats
	log.Println("Testing GetStats...")
	statsResp, err := client.GetStats(ctx, &pb.GetStatsRequest{})
	if err != nil {
		log.Fatalf("GetStats failed: %v", err)
	}
	log.Printf("Stats - URLs: %d, Users: %d\n", statsResp.GetUrls(), statsResp.GetUsers())

	log.Println("All gRPC tests completed successfully!")
}
