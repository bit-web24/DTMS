package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	pb "github.com/bit-web24/DTMS/services/task/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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
	db *gorm.DB
}

func (s *server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	task := &Task{
		ID:          uuid.New().String(),
		Description: req.GetDescription(),
	}

	result := s.db.Create(task)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.CreateTaskResponse{Task: &pb.Task{
		Id:          task.ID,
		Description: task.Description,
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

func (s *server) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	result := s.db.Delete(&Task{}, "id = ?", req.GetId())

	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.DeleteTaskResponse{Success: true}, nil
}

func (s *server) GetAllTasks(ctx context.Context, req *pb.GetAllTasksRequest) (*pb.GetAllTasksResponse, error) {
	var tasks []Task
	result := s.db.Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}

	var pbTasks []*pb.Task
	for _, task := range tasks {
		pbTask := &pb.Task{
			Id:          task.ID,
			Description: task.Description,
		}
		pbTasks = append(pbTasks, pbTask)
	}

	return &pb.GetAllTasksResponse{Tasks: pbTasks}, nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	fmt.Println("DB_HOST: " + os.Getenv("DB_HOST"))

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
	db.AutoMigrate(&Task{})

	lis, err := net.Listen("tcp", ":"+os.Getenv("RPC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &server{db: db})
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("task_service", grpc_health_v1.HealthCheckResponse_SERVING)

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
