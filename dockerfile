FROM golang:1.15-alpine3.12 AS builder

COPY . /github.com/ChingizAdamov/test_work/
WORKDIR  /github.com/ChingizAdamov/test_work/

RUN go mod download
RUN go build -o ./bin/main cmd/main/main.go

FROM alpine:latest

WORKDIR --from=0 /github.com/ChingizAdamov/test_work/bin/main .

EXPOSE 80

CMD ["./main"]