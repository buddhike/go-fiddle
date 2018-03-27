# GO Fiddle

Golang implementation of a HTTP proxy inspired by Fiddler.

## Usage

Run the proxy:

```go
go run cmd/proxy/proxy.go
```

Issue a request:

```sh
# stubbed request
curl https://www.google.com.au/stub -x http://localhost:8080/ -k

# un-stubbed request
curl https://www.google.com.au/ -x http://localhost:8080/ -k
```
