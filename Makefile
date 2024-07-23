PROTOC = protoc
PROTO_DIR = ./proto
SERVICE_TASK_DIR = ./services/task/proto
SERVICE_USER_DIR = ./services/user/proto
GAPI_DIR = ./proto/google

PROTO_FILES = $(PROTO_DIR)/task.proto $(PROTO_DIR)/user.proto

all: generate

generate: $(PROTO_FILES)
	$(PROTOC) -I $(PROTO_DIR) -I $(GAPI_DIR) \
		--go_out=$(SERVICE_TASK_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(SERVICE_TASK_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(SERVICE_TASK_DIR) --grpc-gateway_opt=paths=source_relative \
		$(PROTO_DIR)/task.proto
	$(PROTOC) -I $(PROTO_DIR) -I $(GAPI_DIR) \
		--go_out=$(SERVICE_USER_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(SERVICE_USER_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(SERVICE_USER_DIR) --grpc-gateway_opt=paths=source_relative \
		$(PROTO_DIR)/user.proto

.PHONY: all generate
