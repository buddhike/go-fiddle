package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).
		HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest(shouldInterceptRequest()).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			log.Print(r.URL.String())

			// get stubbed response (a nil response indicates that request should not be stubbed and response should come from actual source)
			return r, stubResponse(r)
		})

	proxy.OnResponse(shouldInterceptResponse()).DoFunc(
		func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			buf, _ := ioutil.ReadAll(r.Body)
			responseStream := ioutil.NopCloser(bytes.NewBuffer(buf))

			s := string(buf)
			log.Print(s)

			r.Body = responseStream

			return r
		})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	log.Print(http.ListenAndServe(fmt.Sprintf(":%s", port), proxy))
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
