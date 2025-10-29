FROM golang:1.24-alpine AS builder

# Устанавливаем нужные пакеты для работы и сборки
RUN apk add --no-cache git tzdata ca-certificates

WORKDIR /app

# Копируем только go модуль и зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем оставшуюся часть проекта
COPY . .


RUN /bin/sh -c 'set -e; \
  for i in 1 2 3; do \
    go install github.com/swaggo/swag/cmd/swag@latest && break || { echo "retry swag ($i)"; sleep 3; }; \
  done && \
  $(go env GOPATH)/bin/swag init'
# Скачиваем swag только в случае необходимости
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    $(go env GOPATH)/bin/swag init

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/server ./main.go

# РUNTIME STAGE
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Копируем собранный бинарник и Swagger docs
COPY --from=builder /app/server /app/server
COPY --from=builder /app/docs /app/docs

ENV GIN_MODE=release
EXPOSE 8000

ENTRYPOINT ["/app/server"]
