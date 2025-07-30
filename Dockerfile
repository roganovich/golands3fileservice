# Берем базовый образ Go
FROM golang:1.18-alpine AS builder
# Устанавливаем зависимости (если нужны)
RUN apk add --no-cache git
# Рабочая директория
WORKDIR /app
# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./
# Скачиваем зависимости
RUN go mod download
# Копируем исходный код
COPY . .
# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o fileservice .
# Создаем финальный образ
FROM alpine:latest
WORKDIR /root/
# Копируем бинарник из builder
COPY --from=builder /app/fileservice .
# Экспонируем порт
EXPOSE 8080
# Команда для запуска
CMD ["./fileservice"]