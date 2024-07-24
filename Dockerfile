FROM golang:1.22.5

WORKDIR /github.com/bit-web24/DTMS/

COPY go.mod .

RUN go mod download
RUN go mod tidy

COPY . .

RUN go build -o dtms main.go
CMD ["./dtms"]
