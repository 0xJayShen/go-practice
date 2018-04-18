FROM scratch

WORKDIR $GOPATH/src/Seckill
COPY . $GOPATH/src/Seckill

EXPOSE 8000
CMD ["./go-gin-example"]