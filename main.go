package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bit-web24/DTMS/health"
	taskpb "github.com/bit-web24/DTMS/services/task/proto"
	userpb "github.com/bit-web24/DTMS/services/user/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	userAddr := os.Getenv("SERVICE_USER_ADDR")
	taskAddr := os.Getenv("SERVICE_TASK_ADDR")

	if userAddr == "" || taskAddr == "" {
		log.Fatalf("Environment variables SERVICE_USER_ADDR or SERVICE_TASK_ADDR are not set")
	}

	userServiceAddr := fmt.Sprintf("%s:50051", userAddr)
	taskServiceAddr := fmt.Sprintf("%s:50052", taskAddr)

	fmt.Println("User Service Address:", userServiceAddr)
	health.CheckHealth(userServiceAddr, "user_service")

	fmt.Println("Task Service Address:", taskServiceAddr)
	health.CheckHealth(taskServiceAddr, "task_service")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, userServiceAddr, opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP gateway for UserService: %v", err)
	}

	err = taskpb.RegisterTaskServiceHandlerFromEndpoint(ctx, mux, taskServiceAddr, opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP gateway for TaskService: %v", err)
	}

	log.Printf("HTTP Gateway listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
