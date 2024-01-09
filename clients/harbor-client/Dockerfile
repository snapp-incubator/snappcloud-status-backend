FROM golang:alpine AS builder

WORKDIR /go/src/app

COPY . .

RUN go build -o app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/src/app/app .

ENV IMAGE_REFERENCE="The image reference URL starts with //."

EXPOSE 8080

CMD ["./app"]
