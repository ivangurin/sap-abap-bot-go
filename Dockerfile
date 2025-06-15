FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./build/bot ./cmd/bot

FROM alpine
WORKDIR /app
COPY --from=builder /app/build/bot ./bot
CMD ["/app/bot"]
