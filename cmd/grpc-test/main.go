package main

import (
	"context"
	"log"
	"time"

	pb "github.com/VladSnap/shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() { // Connect to gRPC server
	conn, err := grpc.NewClient("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewShortenerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test Ping
	log.Println("Testing Ping...")
	pingResp, err := client.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		log.Fatalf("Ping failed: %v", err)
	}
	log.Printf("Ping response: %s\n", pingResp.Status)

	// Test CreateShortLink
	log.Println("Testing CreateShortLink...")
	createResp, err := client.CreateShortLink(ctx, &pb.CreateShortLinkRequest{
		OriginalUrl: "https://example.com",
		UserId:      "test-user-123",
	})
	if err != nil {
		log.Fatalf("CreateShortLink failed: %v", err)
	}
	log.Printf("Created short link: %s (duplicate: %t)\n", createResp.ShortUrl, createResp.IsDuplicate)

	// Test GetStats
	log.Println("Testing GetStats...")
	statsResp, err := client.GetStats(ctx, &pb.GetStatsRequest{})
	if err != nil {
		log.Fatalf("GetStats failed: %v", err)
	}
	log.Printf("Stats - URLs: %d, Users: %d\n", statsResp.Urls, statsResp.Users)

	log.Println("All gRPC tests completed successfully!")
}
