package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	pb "github.com/VladSnap/shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("=== Comprehensive HTTP vs gRPC Comparison Test ===")

	// Test HTTP API
	log.Println("\n--- Testing HTTP API ---")
	testHTTPAPI()

	// Test gRPC API
	log.Println("\n--- Testing gRPC API ---")
	testGRPCAPI()

	log.Println("\n=== All tests completed! ===")
}

func testHTTPAPI() {
	baseURL := "http://localhost:8080"

	// Test 1: Create short link
	log.Println("1. Creating short link via HTTP...")
	resp, err := http.Post(baseURL+"/", "text/plain", bytes.NewBufferString("https://httpbin.org/get"))
	if err != nil {
		log.Fatalf("HTTP POST failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	shortURL := string(body)
	log.Printf("   HTTP short URL created: %s (status: %d)", shortURL, resp.StatusCode)

	// Test 2: Create batch
	log.Println("2. Creating batch via HTTP...")
	batchReq := []map[string]string{
		{"correlation_id": "1", "original_url": "https://httpbin.org/json"},
		{"correlation_id": "2", "original_url": "https://httpbin.org/uuid"},
	}
	batchJSON, _ := json.Marshal(batchReq)

	resp, err = http.Post(baseURL+"/api/shorten/batch", "application/json", bytes.NewBuffer(batchJSON))
	if err != nil {
		log.Fatalf("HTTP batch failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	log.Printf("   HTTP batch response: %s (status: %d)", string(body), resp.StatusCode)
}

func testGRPCAPI() {
	// Connect to gRPC server
	conn, err := grpc.NewClient("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewShortenerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// For testing, let the server generate user IDs automatically
	// In a real client, you would include properly signed auth-cookie metadata

	// Test 1: Create short link
	log.Println("1. Creating short link via gRPC...")
	createResp, err := client.CreateShortLink(ctx, &pb.CreateShortLinkRequest{
		OriginalUrl: "https://httpbin.org/get",
	})
	if err != nil {
		log.Fatalf("gRPC CreateShortLink failed: %v", err)
	}
	log.Printf("   gRPC short URL created: %s (duplicate: %t)", createResp.ShortUrl, createResp.IsDuplicate)

	// Test 2: Create batch
	log.Println("2. Creating batch via gRPC...")
	batchResp, err := client.CreateShortLinkBatch(ctx, &pb.CreateShortLinkBatchRequest{
		Links: []*pb.OriginalLinkBatch{
			{CorrelationId: "1", OriginalUrl: "https://httpbin.org/json"},
			{CorrelationId: "2", OriginalUrl: "https://httpbin.org/uuid"},
		},
	})
	if err != nil {
		log.Fatalf("gRPC CreateShortLinkBatch failed: %v", err)
	}
	log.Printf("   gRPC batch created %d links:", len(batchResp.Links))
	for _, link := range batchResp.Links {
		log.Printf("     - ID %s: %s", link.CorrelationId, link.ShortUrl)
	}

	// Test 3: Get stats
	log.Println("3. Getting stats via gRPC...")
	statsResp, err := client.GetStats(ctx, &pb.GetStatsRequest{})
	if err != nil {
		log.Fatalf("gRPC GetStats failed: %v", err)
	}
	log.Printf("   gRPC stats - URLs: %d, Users: %d", statsResp.Urls, statsResp.Users)

	// Test 4: Ping
	log.Println("4. Ping via gRPC...")
	pingResp, err := client.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		log.Fatalf("gRPC Ping failed: %v", err)
	}
	log.Printf("   gRPC ping response: %s", pingResp.Status)
}
