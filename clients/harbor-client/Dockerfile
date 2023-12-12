FROM golang:1.16-alpine

WORKDIR /client

COPY . .

RUN go mod download

RUN go build -o harbor-client .

ENV IMAGE_REFERENCE "//core.harbor.domain/test-project/busybox:test-tag" 

EXPOSE 8080

CMD ["./harbor-client"]
