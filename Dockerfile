# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY apiApp /app/apiApp
COPY .env /app/.env

WORKDIR /app

CMD ["/app/apiApp"]