# Этап сборки
FROM golang:alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальную часть кода
COPY . .

# Компилируем Go приложение
RUN go build -o server cmd/main.go

# Финальный образ
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add sqlite sqlite-dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированное приложение из этапа сборки
COPY --from=builder /app/server .

# Указываем переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=your_password

# Открываем порт
EXPOSE 7540

# Запускаем сервер
CMD ["./server"]
