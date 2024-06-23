# build a tiny docker image
FROM golang:1.22.2

RUN mkdir /app

COPY go.mod go.sum ./
RUN go mod download

COPY apiApp /app/apiApp
COPY .env /app/.env

WORKDIR /app

CMD ["/app/apiApp"]