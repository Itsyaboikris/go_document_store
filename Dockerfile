FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o docstore ./cmd/main.go

EXPOSE 9000
CMD ["./docstore"]