FROM golang:alpine AS builder

WORKDIR /go/src/app

COPY . .

RUN go build -o app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/src/app/app .

ENV PROXY_URL="The proxy URL goes here."

EXPOSE 8080

CMD ["./app"]
