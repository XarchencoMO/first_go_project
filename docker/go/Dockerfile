# Указываем базовый образ с Go
FROM golang:1.22-alpine

# Создаем рабочую директорию

# Устанавливаем bash и git для работы с Air
#RUN apk add --no-cache bash git

# Устанавливаем Air для горячей перезагрузки
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Компилируем приложение
#RUN go build -o main .

# Открываем порт для API
EXPOSE 8080

# Команда по умолчанию для запуска вашего приложения
#CMD ["./main"]
# Команда по умолчанию для запуска вашего приложения с использованием Air
CMD ["air"]