FROM golang:1.19-alpine3.16 AS service_builder

WORKDIR /

# встановлення додаткових інструментів та бібліотек
RUN apk add gcc libc-dev librdkafka-dev

# встановлення залежностей
COPY go.mod go.sum ./
RUN go mod download

# копіювання основного коду сервісу
COPY services/goods services/goods

# збарання сервісу
WORKDIR /services/goods
RUN go build -ldflags "-w -s -linkmode external -extldflags -static" -a main.go

# підготовка фінального образу
FROM scratch
EXPOSE 8080
COPY --from=service_builder /services/goods/main .
CMD ["./main"]
