FROM golang:1.26.1-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -o /app/exe /app/cmd/todoapp

FROM alpine:3.23
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /app/exe /app/

USER appuser
EXPOSE 5050
CMD ["/app/exe"]
