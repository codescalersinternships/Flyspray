FROM golang:1.20
WORKDIR /app
COPY go.mod go.sum ./ 
RUN go mod download 
COPY . ./
RUN CGO_ENABLED=1 go build -o server ./cmd/server.go
EXPOSE 3000
ENTRYPOINT ["./server"]
CMD ["-p","3000"]