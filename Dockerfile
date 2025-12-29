FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gdtq main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/gdtq .

EXPOSE 8080

CMD ["./gdtq"]