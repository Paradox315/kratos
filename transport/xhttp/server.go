package xhttp

import (
	"context"
	"crypto/tls"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/internal/endpoint"
	"github.com/go-kratos/kratos/v2/internal/host"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"net"
	"net/url"
	"time"
)

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

// Network with server network.
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// Logger with server logger.
func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}

// Listener with server lis
func Listener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

// Middleware with server middleware
func Middleware(ms ...fiber.Handler) ServerOption {
	return func(s *Server) {
		s.ms = ms
	}
}

// Router with server router
func Router(r ...initRouters) ServerOption {
	return func(s *Server) {
		s.router = r
	}
}

// FiberConfig with server config
func FiberConfig(cfg fiber.Config) ServerOption {
	return func(s *Server) {
		s.config = cfg
	}
}

// initRouters is a function to initialize routers.
type initRouters func(r fiber.Router)

type Server struct {
	server   *fiber.App
	baseCtx  context.Context
	lis      net.Listener
	tlsConf  *tls.Config
	endpoint *url.URL
	err      error
	network  string
	address  string
	config   fiber.Config
	ms       []fiber.Handler
	router   []initRouters
	timeout  time.Duration
	log      *log.Helper
}

// NewServer creates an HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		log:     log.NewHelper(log.DefaultLogger),
	}
	for _, o := range opts {
		o(srv)
	}
	srv.server = fiber.New(srv.config)
	for _, m := range srv.ms {
		srv.server.Use(m)
	}
	for _, r := range srv.router {
		r(srv.server)
	}
	srv.err = srv.listenAndEndpoint()
	return srv
}

// Serve serves the server by options.
func (s *Server) Serve() error {
	return s.server.Listener(s.lis)
}

// ServeTLS TODO: not implemented yet
func (s *Server) ServeTLS() error {
	return s.server.Listener(s.lis)
}

// Endpoint returns the endpoint of the server.
func (s *Server) Endpoint() (*url.URL, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

// Start start the FIBER server.
func (s *Server) Start(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}
	s.baseCtx = ctx
	s.log.Infof("[FIBER] server listening on: %s", s.lis.Addr().String())
	var err error
	if s.tlsConf != nil {
		err = s.ServeTLS()
	} else {
		err = s.Serve()
	}
	if !errors.Is(err, fasthttp.ErrConnectionClosed) {
		return err
	}
	return nil
}

// Stop stop the FIBER server.
func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("[FIBER] server stopping")
	return s.server.Shutdown()
}

// Route add a route to the FIBER server.
func (s *Server) Route(init initRouters) {
	init(s.server)
}

// listenAndEndpoint listen and get the endpoint.
func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	s.endpoint = endpoint.NewEndpoint("http", addr, s.tlsConf != nil)
	return nil
}
