FROM golang:latest
LABEL maintainer="yanglei@ks.c.titech.ac.jp"

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]