FROM golang:1.22.7-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o doctorgo ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/doctorgo /app/doctorgo
COPY --from=builder /app/migrations /app/migrations
EXPOSE 8080
CMD ["/app/doctorgo"]
