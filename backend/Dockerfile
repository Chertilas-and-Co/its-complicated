FROM golang:latest AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server/main.go

# Финальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Рабочий каталог
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/static static

EXPOSE 8080
CMD ["/app/server"]
