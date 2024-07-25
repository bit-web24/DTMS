package health

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func CheckHealth(addr string, service string) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &grpc_health_v1.HealthCheckRequest{Service: service}
	res, err := client.Check(ctx, req)
	if err != nil {
		log.Fatalf("could not check health: %v", err)
	}

	fmt.Printf("Health check status: %s\n", res.Status)
}
