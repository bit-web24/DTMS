FROM golang:1.22.5

WORKDIR /github.com/bit-web24/DTMS

COPY go.mod go.sum ./
RUN go mod download

COPY services/task/ services/task/

COPY services/user/proto/ services/user/proto/

RUN go mod tidy

WORKDIR /github.com/bit-web24/DTMS/services/task

RUN go build -o task-service main.go
CMD ["./task-service"]
