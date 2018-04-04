package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"time"

	"go-fiddle/internal/config"
	"go-fiddle/internal/kafkaserver"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elazarl/goproxy"
	"github.com/satori/go.uuid"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	kafkaProducer := kafkaserver.NewProducer()
	requestMap := make(map[*http.Request]string)

	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).
		HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest(shouldInterceptRequest()).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			topic := "request"
			url := r.URL.String()
			log.Print(url)

			httpRequest, _ := httputil.DumpRequest(r, true)

			requestID, _ := uuid.NewV4()
			requestMap[r] = requestID.String()

			prefix := []byte{}
			prefix = append(prefix, []byte(fmt.Sprintf("request-id: %s\r\n", requestID))...)
			prefix = append(prefix, []byte(fmt.Sprintf("timestamp: %s\r\n", time.Now().Format(time.RFC3339)))...)

			kafkaProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          append(prefix, httpRequest...),
			}, nil)

			// get stubbed response (a nil response indicates that request should not be stubbed and response should come from actual source)
			return r, stubResponse(r)
		})

	proxy.OnResponse(shouldInterceptResponse()).DoFunc(
		func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			httpResponse, _ := httputil.DumpResponse(r, true)

			topic := "response"

			requestID := requestMap[r.Request]

			prefix := []byte{}
			prefix = append(prefix, []byte(fmt.Sprintf("request-id: %s\r\n", requestID))...)
			prefix = append(prefix, []byte(fmt.Sprintf("timestamp: %s\r\n", time.Now().Format(time.RFC3339)))...)

			delete(requestMap, r.Request)

			kafkaProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          append(prefix, httpResponse...),
			}, nil)

			return r
		})

	port := config.Get("PORT", "8080")
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), proxy))
}

func shouldInterceptRequest() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		// TODO: query config for whether or not request should be intercepted and logged
		return true
	}
}

func shouldInterceptResponse() goproxy.RespConditionFunc {
	return func(res *http.Response, ctx *goproxy.ProxyCtx) bool {
		// TODO: query config for whether or not request should be intercepted and logged
		return true
	}
}

func stubResponse(req *http.Request) *http.Response {
	// TODO: load stubbing rules from configuration
	if regexp.MustCompile("stub").MatchString(req.RequestURI) {
		return goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusOK, "Stubbed")
	}
	return nil
}
