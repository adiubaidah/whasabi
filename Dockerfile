FROM golang:1.23-alpine as builder

LABEL author="Ahmad Adi Iskandar Ubaidah"


WORKDIR /app
COPY . .


RUN go build -o /app/main


FROM alpine:latest
ENV APP_PORT=5000
EXPOSE ${APP_PORT}
WORKDIR /app
COPY --from=builder /app/main ./
CMD ["./main"]

