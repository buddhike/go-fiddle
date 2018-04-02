FROM go-fiddle-base

WORKDIR /go/src/go-fiddle

RUN go get -u github.com/gorilla/mux

COPY ./ ./

WORKDIR /go/src/go-fiddle/cmd/rest-api
# RUN CGO_ENABLED=0 GOOS=linux go build

ENV PORT=8000
EXPOSE 8000

ENTRYPOINT [ "go", "run", "./main.go" ]
