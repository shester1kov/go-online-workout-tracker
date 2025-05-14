FROM golang:1.24.2-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .
# отключаем cgo и компилируем под Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/app

# FROM alpine:3.21 AS production

# # SSL и таймзоны 
# # RUN apk add --no-cache ca-certificates tzdata

# WORKDIR /root/

# COPY --from=builder /app/server .
# COPY --from=builder /app/migrations ./migrations

CMD ["./server"]