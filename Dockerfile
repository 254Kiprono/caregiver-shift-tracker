# -------- Stage 1: Build --------
    FROM golang:1.23-alpine AS builder

    RUN apk add --no-cache git
    
    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o caregiver-service
    
    # -------- Stage 2: Run --------
    FROM alpine:latest
    
    RUN apk --no-cache add ca-certificates
    
    WORKDIR /root/
    
    COPY --from=builder /app/caregiver-service .
    
    EXPOSE 6000
    
    CMD ["./caregiver-service"]
    