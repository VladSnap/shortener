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

const (
	// Timeout constants.
	defaultTimeout    = 10 * time.Second
	grpcServerAddr    = "127.0.0.1:9090"
	httpServerBaseURL = "http://localhost:8080"
	contentTypePlain  = "text/plain"
	contentTypeJSON   = "application/json"
	correlationID1    = "1"
	correlationID2    = "2"
	testURL1          = "https://httpbin.org/get"
	testURL2          = "https://httpbin.org/json"
	testURL3          = "https://httpbin.org/uuid"
)

// closeBody safely closes response body and logs any errors.
func closeBody(body io.Closer) {
	if err := body.Close(); err != nil {
		log.Printf("Error closing response body: %v", err)
	}
}

// closeConn safely closes connection and logs any errors.
func closeConn(conn io.Closer) {
	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}

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
	// Test 1: Create short link
	log.Println("1. Creating short link via HTTP...")
	resp, err := http.Post(httpServerBaseURL+"/", contentTypePlain, bytes.NewBufferString(testURL1))
	if err != nil {
		log.Fatalf("HTTP POST failed: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	closeBody(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	shortURL := string(body)
	log.Printf("   HTTP short URL created: %s (status: %d)", shortURL, resp.StatusCode)

	// Test 2: Create batch
	log.Println("2. Creating batch via HTTP...")
	batchReq := []map[string]string{
		{"correlation_id": correlationID1, "original_url": testURL2},
		{"correlation_id": correlationID2, "original_url": testURL3},
	}
	batchJSON, err := json.Marshal(batchReq)
	if err != nil {
		log.Fatalf("Failed to marshal batch request: %v", err)
	}

	resp, err = http.Post(httpServerBaseURL+"/api/shorten/batch", contentTypeJSON, bytes.NewBuffer(batchJSON))
	if err != nil {
		log.Fatalf("HTTP batch failed: %v", err)
	}

	body, err = io.ReadAll(resp.Body)
	closeBody(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read batch response body: %v", err)
	}
	log.Printf("   HTTP batch response: %s (status: %d)", string(body), resp.StatusCode)
}

func testGRPCAPI() {
	// Connect to gRPC server
	conn, err := grpc.NewClient(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	client := pb.NewShortenerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)

	// For testing, let the server generate user IDs automatically
	// In a real client, you would include properly signed auth-cookie metadata

	// Test 1: Create short link
	log.Println("1. Creating short link via gRPC...")
	createResp, err := client.CreateShortLink(ctx, &pb.CreateShortLinkRequest{
		OriginalUrl: testURL1,
	})
	if err != nil {
		cancel()
		closeConn(conn)
		log.Fatalf("gRPC CreateShortLink failed: %v", err)
	}
	log.Printf("   gRPC short URL created: %s (duplicate: %t)", createResp.GetShortUrl(), createResp.GetIsDuplicate())

	// Test 2: Create batch
	log.Println("2. Creating batch via gRPC...")
	batchResp, err := client.CreateShortLinkBatch(ctx, &pb.CreateShortLinkBatchRequest{
		Links: []*pb.OriginalLinkBatch{
			{CorrelationId: correlationID1, OriginalUrl: testURL2},
			{CorrelationId: correlationID2, OriginalUrl: testURL3},
		},
	})
	if err != nil {
		log.Fatalf("gRPC CreateShortLinkBatch failed: %v", err)
	}
	log.Printf("   gRPC batch created %d links:", len(batchResp.GetLinks()))
	for _, link := range batchResp.GetLinks() {
		log.Printf("     - ID %s: %s", link.GetCorrelationId(), link.GetShortUrl())
	}

	// Test 3: Get stats
	log.Println("3. Getting stats via gRPC...")
	statsResp, err := client.GetStats(ctx, &pb.GetStatsRequest{})
	if err != nil {
		log.Fatalf("gRPC GetStats failed: %v", err)
	}
	log.Printf("   gRPC stats - URLs: %d, Users: %d", statsResp.GetUrls(), statsResp.GetUsers())

	// Test 4: Ping
	log.Println("4. Ping via gRPC...")
	pingResp, err := client.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		log.Fatalf("gRPC Ping failed: %v", err)
	}
	log.Printf("   gRPC ping response: %s", pingResp.GetStatus())

	// Clean up resources
	cancel()
	closeConn(conn)
}
