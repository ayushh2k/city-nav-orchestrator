FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./bin/orchestrator ./cmd/orchestrator/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/orchestrator .

RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["./orchestrator"]