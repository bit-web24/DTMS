all: target

target:
	protoc --go_out=. --go-grpc_out=. proto/user.proto proto/task.proto
