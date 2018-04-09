package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"time"

	"go-fiddle/cmd/config"
	"go-fiddle/cmd/internal/database"
	"go-fiddle/cmd/internal/kafkaserver"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elazarl/goproxy"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	kafkaProducer := kafkaserver.NewProducer()
	requestMap := make(map[*http.Request]string)

	session := database.GetDatabaseConnection()
	defer session.Close()
	// session.SetMode(mgo.Monotonic, true)
	collection := database.GetDatabaseCollection(session, "messages")

	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).
		HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest(shouldInterceptRequest()).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			topic := "request"
			url := r.URL.String()
			log.Print(url)

			httpMessage := HTTPMessage{}
			request, _ := httputil.DumpRequest(r, true)
			requestID, _ := uuid.NewV4()
			requestMap[r] = requestID.String()
			timestamp := time.Now()

			httpRequest := unmarshalHTTPRequest(request)
			httpRequest.Timestamp = &timestamp

			httpMessage.ID = requestID.String()
			httpMessage.Request = httpRequest

			prefix := []byte{}
			prefix = append(prefix, []byte(fmt.Sprintf("request-id: %s\r\n", requestID))...)
			prefix = append(prefix, []byte(fmt.Sprintf("timestamp: %s\r\n", time.Now().Format(time.RFC3339)))...)

			go func() {
				err := collection.Insert(httpMessage)

				if err != nil {
					log.Fatal(err)
				}
			}()

			if jsonMessage, err := json.Marshal(summariseMessage(httpMessage)); err == nil {
				kafkaProducer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Value:          jsonMessage,
				}, nil)
			}

			// get stubbed response (a nil response indicates that request should not be stubbed and response should come from actual source)
			return r, stubResponse(r)
		})

	proxy.OnResponse(shouldInterceptResponse()).DoFunc(
		func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			httpResponse, _ := httputil.DumpResponse(r, false)
			buf, _ := ioutil.ReadAll(r.Body)
			responseStream := ioutil.NopCloser(bytes.NewBuffer(buf))
			httpResponse = append(httpResponse, buf...)
			httpMessage := HTTPMessage{}

			r.Body = responseStream

			timestamp := time.Now()
			topic := "response"

			requestID := requestMap[r.Request]
			err := collection.FindId(requestID).One(&httpMessage)

			if err != nil {
				log.Fatal(err)
			}

			response := unmarshalHTTPResponse(httpResponse)
			response.Timestamp = &timestamp
			httpMessage.Response = response

			delete(requestMap, r.Request)

			go func() {
				err := collection.Update(bson.M{"_id": requestID}, httpMessage)

				if err != nil {
					log.Fatal(err)
				}
			}()

			if jsonMessage, err := json.Marshal(summariseMessage(httpMessage)); err == nil {
				kafkaProducer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Value:          jsonMessage,
				}, nil)
			}

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
