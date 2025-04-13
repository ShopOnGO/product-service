FROM golang:1.23.3 AS builder

WORKDIR /product

# Устанавливаем pg_isready и очищаем кеш
RUN apt-get update && apt-get install -y postgresql-client \
    && rm -rf /var/lib/apt/lists/* && apt-get clean

# Отключаем CGO для статической компиляции
 ENV CGO_ENABLED=0

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download && go mod verify

# Копируем весь код
COPY . .

# Компилируем бинарник
RUN go build -o /product/product_service ./cmd/server.go



# Второй этап: финальный образ (без лишних инструментов)
FROM alpine:latest

WORKDIR /product

# Устанавливаем postgresql-client и dos2unix
RUN apk add --no-cache postgresql-client dos2unix

COPY .env /product/.env

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /product/product_service /product/product_service

# Копируем wait-for-db.sh и делаем исполняемым
COPY --from=builder /product/wait-for-db.sh /product/wait-for-db.sh
RUN chmod +x /product/wait-for-db.sh

# Преобразуем формат строки в скрипте wait-for-db.sh в Unix-формат
RUN dos2unix /product/wait-for-db.sh

# Запуск приложения
CMD ["/product/product_service"]
