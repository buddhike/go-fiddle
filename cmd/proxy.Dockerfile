FROM go-fiddle-base

WORKDIR /go/src/go-fiddle/cmd

RUN go get -u github.com/elazarl/goproxy
RUN go get -u github.com/satori/go.uuid

COPY internal internal
COPY config config
COPY proxy proxy

WORKDIR /go/src/go-fiddle/cmd/proxy
RUN CGO_ENABLED=1 GOOS=linux go build

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT [ "./proxy" ]
