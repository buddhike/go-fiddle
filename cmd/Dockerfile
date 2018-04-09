FROM golang:1.10.0-alpine3.7

RUN apk add --update --no-cache \
      alpine-sdk \
      bash \
      ca-certificates \
      libressl \
      tar \
      git \
      openssh \
      openssl \
      yajl-dev \
      zlib-dev \
      cyrus-sasl-dev \
      openssl-dev \
      build-base \
      coreutils

WORKDIR /root
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN /root/librdkafka/configure
RUN make
RUN make install
#For golang applications
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN apk add --no-cache librdkafka-dev git gcc

RUN go get -u github.com/confluentinc/confluent-kafka-go/kafka

WORKDIR /go/src/go-fiddle/cmd

RUN go get -u github.com/elazarl/goproxy
RUN go get -u github.com/satori/go.uuid
RUN go get -u gopkg.in/mgo.v2

COPY internal internal
COPY config config
COPY proxy proxy

WORKDIR /go/src/go-fiddle/cmd/proxy
RUN CGO_ENABLED=1 GOOS=linux go build

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT [ "./proxy" ]
