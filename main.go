package main

import (
	"context"
	"log"
	"net/http"

	taskpb "github.com/bit-web24/DTMS/services/task/proto"
	userpb "github.com/bit-web24/DTMS/services/user/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register UserService handler from endpoint
	// user_service_addr := fmt.Sprintf("%s:50051", os.Getenv("SERVICE_USER_ADDR"))
	// fmt.Printf("TaskServiceAddr: %s\n", os.Getenv("SERVICE_USER_ADDR"))
	err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "0.0.0.0:50051", opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP gateway for UserService: %v", err)
	}

	// Register TaskService handler from endpoint
	// task_service_addr := fmt.Sprintf("%s:50052", os.Getenv("SERVICE_TASK_ADDR"))
	// fmt.Printf("UserServiceAddr: %s\n", os.Getenv("SERVICE_TASK_ADDR"))
	err = taskpb.RegisterTaskServiceHandlerFromEndpoint(ctx, mux, "0.0.0.0:50052", opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP gateway for TaskService: %v", err)
	}

	log.Printf("HTTP Gateway listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
