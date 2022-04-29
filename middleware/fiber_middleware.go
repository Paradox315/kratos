package middleware

import (
	"github.com/go-kratos/kratos/v2/transport/xhttp/apistate"
	"github.com/gofiber/fiber/v2"
)

// Constants for the middleware
const (
	SupportPackageIsVersion1 = true
	AuthenticatorCfg         = "Authenticator"
	AuthorizerCfg            = "Authorizer"
	OperationsCfg            = "Operations"
	CacheCfg                 = "Cache"
	LimiterCfg               = "Limiter"
)

var middlewareConf = map[string]FiberMiddleware{
	AuthenticatorCfg: defaultMiddleware(),
	AuthorizerCfg:    defaultMiddleware(),
	OperationsCfg:    defaultMiddleware(),
	CacheCfg:         defaultMiddleware(),
	LimiterCfg:       defaultMiddleware(),
}

// FiberMiddleware is a middleware for Fiber
type FiberMiddleware interface {
	MiddlewareFunc() fiber.Handler
	Name() string
}

// UnimplementedMiddleware is a middleware that is not implemented
type UnimplementedMiddleware struct {
}

// MiddlewareFunc returns the middleware function
func (u *UnimplementedMiddleware) MiddlewareFunc() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return apistate.Error[any]().WithMessage("Middleware not implemented").Send(c)
	}
}

// Name returns the name of the middleware
func (u *UnimplementedMiddleware) Name() string {
	return "UnimplementedMiddleware"
}

// defaultMiddleware returns the default unimplemented middleware
func defaultMiddleware() *UnimplementedMiddleware {
	return &UnimplementedMiddleware{}
}

// RegisterMiddleware registers a middleware
func RegisterMiddleware(mw FiberMiddleware) {
	if middlewareConf == nil {
		middlewareConf = make(map[string]FiberMiddleware)
	}
	middlewareConf[mw.Name()] = mw
}

// Authenticator returns the Authenticator middleware
func Authenticator() fiber.Handler {
	if mw, ok := middlewareConf[AuthenticatorCfg]; ok {
		return mw.MiddlewareFunc()
	}
	return middlewareConf[AuthenticatorCfg].MiddlewareFunc()
}

// Authorizer returns the Authorizer middleware
func Authorizer() fiber.Handler {
	if mw, ok := middlewareConf[AuthorizerCfg]; ok {
		return mw.MiddlewareFunc()
	}
	return middlewareConf[AuthorizerCfg].MiddlewareFunc()
}

// Cache returns the Cache middleware
func Cache() fiber.Handler {
	if mw, ok := middlewareConf[CacheCfg]; ok {
		return mw.MiddlewareFunc()
	}
	return middlewareConf[CacheCfg].MiddlewareFunc()
}

// Limiter returns the Limiter middleware
func Limiter() fiber.Handler {
	if mw, ok := middlewareConf[LimiterCfg]; ok {
		return mw.MiddlewareFunc()
	}
	return middlewareConf[LimiterCfg].MiddlewareFunc()
}

// Operations returns the Operations middleware
func Operations() fiber.Handler {
	if mw, ok := middlewareConf[OperationsCfg]; ok {
		return mw.MiddlewareFunc()
	}
	return middlewareConf[OperationsCfg].MiddlewareFunc()
}

// CustomMiddleware returns a custom middleware with your config key
func CustomMiddleware(name string) fiber.Handler {
	if mw, ok := middlewareConf[name]; ok {
		return mw.MiddlewareFunc()
	}
	middlewareConf[name] = defaultMiddleware()
	return middlewareConf[name].MiddlewareFunc()
}
