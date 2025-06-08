# Build stage
FROM golang:1.21.4-alpine3.18 AS builder
WORKDIR /app

# Копируем исходный код и зависимости
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/mortgage-calculator ./cmd

# Final stage
FROM alpine:3.18
WORKDIR /app

# Копируем бинарник и конфигурацию
COPY --from=builder /app/mortgage-calculator /app/mortgage-calculator
COPY config.yml /app/config.yml

# Устанавливаем права и пользователя
RUN chmod +x /app/mortgage-calculator && \
    adduser -D appuser && \
    chown appuser:appuser /app/mortgage-calculator

USER appuser

EXPOSE 8080
CMD ["./mortgage-calculator"]