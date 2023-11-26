FROM golang:alpaine

WORKDIR /client

COPY . .

RUN go build -o harbor-client

EXPOSE 8080

CMD ["./harbor-client"]
