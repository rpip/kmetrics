FROM golang:1.12.6-alpine3.9

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN set -x \
  && apk add --no-cache --virtual .build-deps git


RUN go build -o main .
CMD ["/app/main"]
