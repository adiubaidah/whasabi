FROM golang:1.23-alpine as builder

LABEL author="Ahmad Adi Iskandar Ubaidah"

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY . .

ENV CGO_ENABLED=1
RUN go build -o /app/main


FROM alpine:latest
ENV APP_PORT=5000
EXPOSE ${APP_PORT}
WORKDIR /app
COPY --from=builder /app/main ./
CMD ["./main"]

