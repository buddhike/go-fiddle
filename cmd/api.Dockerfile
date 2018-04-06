FROM go-fiddle-base

WORKDIR /go/src/go-fiddle/cmd

RUN go get -u github.com/gorilla/handlers
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/gorilla/websocket
RUN go get -u gopkg.in/mgo.v2

COPY config config
COPY internal internal
COPY rest-api rest-api

WORKDIR /go/src/go-fiddle/cmd/rest-api
RUN CGO_ENABLED=1 GOOS=linux go build

ENV PORT=8000
EXPOSE 8000

ENTRYPOINT [ "./rest-api" ]
