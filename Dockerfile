# Используем официальный образ Golang
FROM golang:1.22.4 as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем приложение
RUN go build -o main .

# Используем минимальный образ для запуска
FROM ubuntu:latest

# Копируем исполняемый файл и фронтенд
COPY --from=builder /app/main /app/main
COPY web /app/web

# Устанавливаем переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/db.sqlite
ENV TODO_PASSWORD=your_password

# Открываем порт
EXPOSE 7540

# Команда для запуска приложения
CMD ["/app/main"]
