FROM golang:1.21.4 as builder

# Рабочая директория внутри контейнера
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Копирование кода 
COPY . .

# Сборка 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

CMD ["./main"]