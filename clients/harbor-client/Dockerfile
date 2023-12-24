FROM golang:1.16-alpine

WORKDIR /client

COPY . .

RUN go mod download

RUN go build -o harbor-client .

ENV IMAGE_REFERENCE "empty ref" 

EXPOSE 8080

CMD ["./harbor-client"]
