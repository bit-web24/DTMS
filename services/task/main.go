package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	userpb "github.com/bit-web24/DTMS/services/user/proto"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	pb "github.com/bit-web24/DTMS/services/task/proto"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Task struct {
	ID          string    `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Description string    `gorm:"size:255;not null"`
	UserID      string    `gorm:"size:255"`
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp"`
}

type server struct {
	pb.UnimplementedTaskServiceServer
	db         *gorm.DB
	userClient userpb.UserServiceClient
}

func (s *server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if req.GetUserId() != "" {
		userReq := &userpb.GetUserRequest{Id: req.GetUserId()}
		_, err := s.userClient.GetUser(ctx, userReq)
		if err != nil {
			return nil, fmt.Errorf("user not found")
		}
	}

	task := &Task{
		ID:          uuid.New().String(),
		Description: req.GetDescription(),
		UserID:      req.GetUserId(),
	}
	result := s.db.Create(task)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.CreateTaskResponse{Task: &pb.Task{
		Id:          task.ID,
		Description: task.Description,
		UserId:      task.UserID,
	},
	}, nil
}

func (s *server) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	var task Task
	result := s.db.First(&task, "id = ?", req.GetId())
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.GetTaskResponse{Task: &pb.Task{
		Id:          task.ID,
		Description: task.Description,
		UserId:      task.UserID,
	},
	}, nil
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
		os.Getenv("DB_TASK_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIME_ZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&Task{})

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userConn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &server{db: db, userClient: userClient})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
