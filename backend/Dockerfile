FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o difference-engine .

FROM alpine:3.21

RUN apk add --no-cache ffmpeg

WORKDIR /app
COPY --from=builder /build/difference-engine ./
COPY frontend/ ./frontend/

EXPOSE 8000
CMD ["./difference-engine"]
