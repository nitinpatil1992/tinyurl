FROM  golang:1.11

ENV GO111MODULE=on

WORKDIR /go/src/main
COPY . .

RUN go build -o tinyurl
CMD ["./tinyurl"]