# Build stage
FROM golang:1.24.6-alpine3.22 AS builder
WORKDIR /app

ARG DB_USER
ARG DB_PASSWORD
ARG DB_DRIVER
ARG DB_NAME
ARG DB_SOURCE
ARG DB_SOURCE_DOCKER
ARG SERVER_ADDRESS
ARG TOKEN_SYMMETRIC_KEY
ARG ACCESS_TOKEN_DURATION

ENV DB_USER=$DB_USER
ENV DB_PASSWORD=$DB_PASSWORD
ENV DB_DRIVER=$DB_DRIVER
ENV DB_NAME=$DB_NAME
ENV DB_SOURCE=$DB_SOURCE
ENV SERVER_ADDRESS=$SERVER_ADDRESS
ENV TOKEN_SYMMETRIC_KEY=$TOKEN_SYMMETRIC_KEY
ENV ACCESS_TOKEN_DURATION=$ACCESS_TOKEN_DURATION

COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz


# Run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY start.sh .
COPY wait-for.sh .
RUN chmod +x start.sh wait-for.sh
COPY db/migration ./migration

EXPOSE  8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]