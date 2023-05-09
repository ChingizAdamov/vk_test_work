FROM golang:1.20.4-alpine3.16 AS builder

COPY . /github.com/ChingizAdamov/test_work/
WORKDIR /github.com/ChingizAdamov/test_work/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/ChingizAdamov/test_work/bin/bot .

EXPOSE 80

CMD ["./bot"]