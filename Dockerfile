FROM golang:1.25.0-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o doctor_go ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/doctor_go /app/doctor_go
COPY --from=builder /app/migrations /app/migrations
EXPOSE 8080
CMD ["/app/doctor_go"]
