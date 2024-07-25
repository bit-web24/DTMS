package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	pb "github.com/bit-web24/DTMS/services/user/proto"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        string    `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Username  string    `gorm:"size:255;not null;unique"`
	Email     string    `gorm:"size:255;not null;unique"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}

type server struct {
	pb.UnimplementedUserServiceServer
	db *gorm.DB
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &User{ID: uuid.New().String(), Username: req.GetUsername(), Email: req.GetEmail()}
	result := s.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pb.CreateUserResponse{User: &pb.User{Id: user.ID, Username: user.Username, Email: user.Email}}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var user User
	result := s.db.First(&user, "id = ?", req.GetId())
	if result.Error != nil {
		return nil, result.Error
	}
	return &pb.GetUserResponse{User: &pb.User{Id: user.ID, Username: user.Username, Email: user.Email}}, nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIME_ZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&User{})

	lis, err := net.Listen("tcp", ":"+os.Getenv("RPC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: db})
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("user_service", grpc_health_v1.HealthCheckResponse_SERVING)

	// Start the gRPC server in a separate goroutine
	go func() {
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Start an HTTP server for the health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	httpPort := os.Getenv("HTTP_PORT")
	log.Printf("HTTP server listening on port %s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
