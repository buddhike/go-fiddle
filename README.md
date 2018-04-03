# GO Fiddle

Golang implementation of a HTTP proxy inspired by Fiddler.

## Getting started

Dependencies:

* `golang`
* `kafka`
* `mongodb`

Alternatively:

* `docker`
* `docker-compose`

Cloning:

```sh
# clone into GOPATH
cd $(go env GOPATH)
git clone https://github.com/socsieng/go-fiddle.git
cd go-fiddle
```

## Usage

Using `docker-compose`:

```sh
docker-compose up
```

Issue a request:

```sh
# stubbed request
curl https://www.google.com.au/stub -x http://localhost:8080/ -k

# un-stubbed request
curl https://www.google.com.au/ -x http://localhost:8080/ -k
```
