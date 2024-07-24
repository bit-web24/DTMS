FROM golang:1.22.5

WORKDIR /github.com/bit-web24/DTMS

COPY go.mod go.sum ./
RUN go mod download

COPY services/user/ services/user/

COPY services/task/proto/ services/task/proto/

RUN go mod tidy

WORKDIR /github.com/bit-web24/DTMS/services/user

RUN go build -o user-service main.go
CMD ["./user-service"]
