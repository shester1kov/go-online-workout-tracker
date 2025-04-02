FROM golang:1.23.1-bookworm

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o server ./cmd/app

EXPOSE ${PORT}

CMD ["./server"]