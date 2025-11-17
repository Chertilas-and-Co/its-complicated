FROM golang:latest AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backend .

# Финальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Рабочий каталог
WORKDIR /app
COPY --from=builder /app/backend .
COPY --from=builder /app/static static

EXPOSE 8080
CMD ["/app/backend"]
