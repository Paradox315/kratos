package xhttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/internal/host"
	"github.com/go-kratos/kratos/v2/internal/httputil"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/gofiber/fiber/v2"
	"time"
)

// DecodeErrorFunc is decode error func.
type DecodeErrorFunc func(ctx context.Context, resp *Response) error

// EncodeRequestFunc is request encode func.
type EncodeRequestFunc func(ctx context.Context, contentType string, in interface{}) (body []byte, err error)

// DecodeResponseFunc is response decode func.
type DecodeResponseFunc func(ctx context.Context, resp *Response, out interface{}) error

// callOptions is call options.
type callOption struct {
	path        string
	method      string
	contentType string
}

type Response struct {
	Code        int
	Body        []byte
	Errors      []error
	ContentType string
}

// defaultCallOptions is the default set of call options.
func defaultCallOption(path, method string) callOption {
	return callOption{
		path:        path,
		method:      method,
		contentType: fiber.MIMEApplicationJSON,
	}
}

// SetContentType sets the content type for a call option.
func (o *callOption) SetContentType(contentType string) {
	o.contentType = contentType
}

// ClientOption is FIBER client option.
type ClientOption func(*clientOptions)

// Client is an FIBER transport client.
type clientOptions struct {
	ctx          context.Context
	tlsConf      *tls.Config
	timeout      time.Duration
	endpoint     string
	userAgent    string
	encoder      EncodeRequestFunc
	decoder      DecodeResponseFunc
	errorDecoder DecodeErrorFunc
	selector     selector.Selector
	discovery    registry.Discovery
	block        bool
}

// WithTimeout with client request timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = d
	}
}

// WithRequestEncoder with client request encoder.
func WithRequestEncoder(encoder EncodeRequestFunc) ClientOption {
	return func(o *clientOptions) {
		o.encoder = encoder
	}
}

// WithResponseDecoder with client response decoder.
func WithResponseDecoder(decoder DecodeResponseFunc) ClientOption {
	return func(o *clientOptions) {
		o.decoder = decoder
	}
}

// WithErrorDecoder with client error decoder.
func WithErrorDecoder(errorDecoder DecodeErrorFunc) ClientOption {
	return func(o *clientOptions) {
		o.errorDecoder = errorDecoder
	}
}

// WithUserAgent with client user agent.
func WithUserAgent(ua string) ClientOption {
	return func(o *clientOptions) {
		o.userAgent = ua
	}
}

// WithEndpoint with client addr.
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithDiscovery with client discovery.
func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
	}
}

// WithSelector with client selector.
func WithSelector(selector selector.Selector) ClientOption {
	return func(o *clientOptions) {
		o.selector = selector
	}
}

// WithBlock with client block.
func WithBlock() ClientOption {
	return func(o *clientOptions) {
		o.block = true
	}
}

// WithTLSConfig with tls config.
func WithTLSConfig(c *tls.Config) ClientOption {
	return func(o *clientOptions) {
		o.tlsConf = c
	}
}

// Client is an HTTP client.
type Client struct {
	opts     clientOptions
	target   *Target
	r        *resolver
	cc       *fiber.Agent
	insecure bool
}

// NewClient returns an HTTP client.
func NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	options := clientOptions{
		ctx:          ctx,
		timeout:      2000 * time.Millisecond,
		encoder:      DefaultRequestEncoder,
		decoder:      DefaultResponseDecoder,
		errorDecoder: DefaultErrorDecoder,
		selector:     wrr.New(),
	}
	for _, o := range opts {
		o(&options)
	}

	insecure := options.tlsConf == nil
	target, err := parseTarget(options.endpoint, insecure)
	if err != nil {
		return nil, err
	}
	var r *resolver
	if options.discovery != nil {
		if target.Scheme == "discovery" {
			if r, err = newResolver(ctx, options.discovery, target, options.selector, options.block, insecure); err != nil {
				return nil, fmt.Errorf("[fiber client] new resolver failed!err: %v", options.endpoint)
			}
		} else if _, _, err := host.ExtractHostPort(options.endpoint); err != nil {
			return nil, fmt.Errorf("[fiber client] invalid endpoint format: %v", options.endpoint)
		}
	}
	agent := fiber.AcquireAgent()
	agent.JSONEncoder(encoding.GetCodec("json").Marshal)
	agent.JSONDecoder(encoding.GetCodec("json").Unmarshal)
	agent.Timeout(options.timeout)
	return &Client{
		opts:     options,
		target:   target,
		insecure: insecure,
		r:        r,
		cc:       agent,
	}, nil
}

// Invoke makes an rpc call procedure for remote service.
func (client *Client) Invoke(ctx context.Context, method, path string, args interface{}, reply interface{}) (err error) {
	c := defaultCallOption(path, method)
	agent := client.cc
	var body []byte
	req := agent.Request()
	switch c.method {
	case fiber.MethodGet:
		req.Header.SetMethod(fiber.MethodGet)
	case fiber.MethodPost:
		req.Header.SetMethod(fiber.MethodPost)
	case fiber.MethodPut:
		req.Header.SetMethod(fiber.MethodPut)
	case fiber.MethodDelete:
		req.Header.SetMethod(fiber.MethodDelete)
	case fiber.MethodPatch:
		req.Header.SetMethod(fiber.MethodPatch)
	default:
		return errors.BadRequest("[fiber client] invalid method: %s", c.method)
	}
	if args != nil {
		body, err = client.opts.encoder(ctx, c.contentType, args)
		if err != nil {
			return err
		}
	}
	var url string
	if client.insecure {
		url = fmt.Sprintf("%s://%s%s", "http", client.target.Authority, path)
	} else {
		url = fmt.Sprintf("%s://%s%s", "https", client.target.Authority, path)
	}
	req.SetRequestURI(url)
	agent.Body(body)
	if err != nil {
		return err
	}
	if c.contentType != "" {
		agent.ContentType(c.contentType)
	}
	if client.opts.userAgent != "" {
		agent.UserAgent(client.opts.userAgent)
	}
	return client.invoke(ctx, c, args, reply)
}

func (client *Client) invoke(ctx context.Context, opt callOption, args interface{}, reply interface{}) error {
	h := func(ctx context.Context, in interface{}) (interface{}, error) {
		res, err := client.do(opt)
		if err != nil {
			return nil, err
		}
		if err := client.opts.decoder(ctx, res, reply); err != nil {
			return nil, err
		}
		return reply, nil
	}
	_, err := h(ctx, args)
	return err
}

func (client *Client) do(opt callOption) (*Response, error) {
	var resp *Response
	var done func(context.Context, selector.DoneInfo)
	if client.r != nil {
		var (
			err  error
			node selector.Node
		)
		if node, done, err = client.opts.selector.Select(context.Background()); err != nil {
			return nil, errors.ServiceUnavailable("NODE_NOT_FOUND", err.Error())
		}
		client.cc.Host(node.Address())
	}
	err := client.cc.Parse()
	if err == nil {
		err = client.opts.errorDecoder(context.Background(), resp)
	}

	if err != nil {
		return nil, err
	}
	code, body, errs := client.cc.Bytes()
	resp = &Response{
		Code:        code,
		Body:        body,
		Errors:      errs,
		ContentType: opt.contentType,
	}
	if done != nil {
		done(context.Background(), selector.DoneInfo{Err: err})
	}
	return resp, nil
}

// Close tears down the Transport and all underlying connections.
func (client *Client) Close() error {
	if client.r != nil {
		return client.r.Close()
	}
	return nil
}

// DefaultRequestEncoder is an HTTP request encoder.
func DefaultRequestEncoder(ctx context.Context, contentType string, in interface{}) ([]byte, error) {
	name := httputil.ContentSubtype(contentType)
	body, err := encoding.GetCodec(name).Marshal(in)
	if err != nil {
		return nil, err
	}
	return body, err
}

// DefaultResponseDecoder is an FIBER response decoder.
func DefaultResponseDecoder(ctx context.Context, resp *Response, v interface{}) error {
	return CodecForResponse(resp).Unmarshal(resp.Body, v)
}

// DefaultErrorDecoder is an FIBER error decoder.
func DefaultErrorDecoder(ctx context.Context, resp *Response) (err error) {
	if resp.Code >= 200 && resp.Code <= 299 {
		return nil
	}
	e := new(errors.Error)
	if err = CodecForResponse(resp).Unmarshal(resp.Body, e); err == nil {
		e.Code = int32(resp.Code)
		return e
	}
	return errors.Errorf(resp.Code, errors.UnknownReason, err.Error())
}

// CodecForResponse gets the codec for the response.
func CodecForResponse(resp *Response) encoding.Codec {
	codec := encoding.GetCodec(httputil.ContentSubtype(resp.ContentType))
	if codec != nil {
		return codec
	}
	return encoding.GetCodec("json")
}
