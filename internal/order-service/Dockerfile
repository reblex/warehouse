FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY internal/order-service/*.go ./

RUN go build -o /service

EXPOSE 8000

CMD ["/service"]