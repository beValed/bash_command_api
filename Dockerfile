FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache bash
RUN go mod tidy
RUN go build -o main ./cmd/server

CMD ["./main"]
