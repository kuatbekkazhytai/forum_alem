FROM golang:latest

RUN mkdir -p /usr/src/forum/
WORKDIR /usr/src/forum/

COPY . /usr/src/forum/
RUN go get github.com/mattn/go-sqlite3
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/satori/go.uuid

EXPOSE 8081

RUN go build -o main .
CMD ["/usr/src/forum/main"]