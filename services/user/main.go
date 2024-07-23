package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/bit-web24/DTMS/services/user/proto"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"google.golang.org/grpc"
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
		os.Getenv("DB_USER_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIME_ZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&User{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: db})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
