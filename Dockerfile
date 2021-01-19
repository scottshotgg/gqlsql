FROM golang:1.14-alpine3.12 AS builder

WORKDIR /app

COPY . .

RUN go build -o main

FROM alpine:3.12

COPY --from=builder /app/main /bin/app

CMD ["/bin/app"]