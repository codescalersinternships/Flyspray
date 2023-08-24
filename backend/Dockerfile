FROM golang:1.20.5

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -a -o server ./cmd/server.go

CMD [ "./server", "-f", "./config.json" ]
