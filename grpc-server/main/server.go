package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	pb "github.com/prrrrnav/proto-rep/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type server struct {
	pb.UnimplementedAuditServiceServer
	redisClient *redis.Client
}

// CreateEvent handles incoming gRPC requests and publishes events to Redis
func (s *server) CreateEvent(ctx context.Context, event *pb.AuditEvent) (*pb.EventResponse, error) {
	log.Printf("Received event: %v", event)

	// Convert the event to JSON
	jsonData, err := protojson.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return &pb.EventResponse{Status: "Error", Message: "Failed to marshal event"}, err
	}

	// Publish the event to the Redis Stream
	streamName := "audit_events"
	eventID := time.Now().UnixNano()

	_, err = s.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"event_id": eventID, // Include event ID in the stream
			"body":     string(jsonData), // Store the entire event JSON in the "body" field
		},
	}).Result()

	if err != nil {
		log.Printf("Failed to publish event to Redis: %v", err)
		return &pb.EventResponse{Status: "Error", Message: "Failed to publish event"}, err
	}

	log.Println("Event published to Redis Stream successfully.")
	return &pb.EventResponse{Status: "Success", Message: "Event received and published to Redis"}, nil
}

func main() {
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "0.0.0.0:6379",
	})

	// Test Redis connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuditServiceServer(grpcServer, &server{redisClient: redisClient})

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	log.Println("gRPC server with Redis integration listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
