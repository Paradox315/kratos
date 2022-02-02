package xhttp

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

var _ Transporter = &Transport{}

// Transporter is fasthttp Transporter
type Transporter interface {
	transport.Transporter
	Request() *fasthttp.Request
	PathTemplate() string
}

// Transport is an FASTHTTP transport.
type Transport struct {
	endpoint     string
	operation    string
	reqHeader    reqHeaderCarrier
	replyHeader  replyHeaderCarrier
	request      *fasthttp.Request
	pathTemplate string
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind {
	return transport.KindXHTTP
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

// Request returns the HTTP request.
func (tr *Transport) Request() *fasthttp.Request {
	return tr.request
}

// RequestHeader returns the request header.
func (tr *Transport) RequestHeader() transport.Header {
	return &tr.reqHeader
}

// ReplyHeader returns the reply header.
func (tr *Transport) ReplyHeader() transport.Header {
	return &tr.replyHeader
}

// PathTemplate returns the http path template.
func (tr *Transport) PathTemplate() string {
	return tr.pathTemplate
}

// SetOperation sets the transport operation.
func SetOperation(ctx context.Context, op string) {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if tr, ok := tr.(*Transport); ok {
			tr.operation = op
		}
	}
}

type reqHeaderCarrier struct {
	reqHeader fiber.Ctx
}

// Get returns the value associated with the passed key.
func (hc *reqHeaderCarrier) Get(key string) string {
	return hc.reqHeader.Get(key)
}

// Set stores the key-value pair.
func (hc *reqHeaderCarrier) Set(key string, value string) {
	hc.reqHeader.Request().Header.Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc *reqHeaderCarrier) Keys() []string {
	keys := make([]string, 0, hc.reqHeader.Request().Header.Len())
	for key := range hc.reqHeader.GetReqHeaders() {
		keys = append(keys, key)
	}
	return keys
}

type replyHeaderCarrier struct {
	replyHeader fiber.Ctx
}

// Get returns the value associated with the passed key.
func (hc *replyHeaderCarrier) Get(key string) string {
	return hc.replyHeader.GetRespHeader(key)
}

// Set stores the key-value pair.
func (hc *replyHeaderCarrier) Set(key string, value string) {
	hc.replyHeader.Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc *replyHeaderCarrier) Keys() []string {
	keys := make([]string, 0, hc.replyHeader.Request().Header.Len())
	for key, _ := range hc.replyHeader.GetRespHeaders() {
		keys = append(keys, key)
	}
	return keys
}
