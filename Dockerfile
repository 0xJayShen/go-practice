FROM golang:latest

WORKDIR $GOPATH/src/gin-docker-mysql
COPY . $GOPATH/src/gin-docker-mysql

EXPOSE 8000
RUN go build .

EXPOSE 8000
ENTRYPOINT ["./gin-docker-mysql"]