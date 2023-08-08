FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o server ./cmd/server.go

FROM alpine:3

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 3000

ENTRYPOINT ["./server"]

CMD ["-p", "3000"]