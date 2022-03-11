package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// Constants for the middleware
const (
	SupportPackageIsVersion1 = true
	AuthenticatorCfg         = "Authenticator"
	AuthorizerCfg            = "Authorizer"
	OperationsCfg            = "Operations"
)

var middlewareConf = map[string]FiberMiddleware{
	AuthenticatorCfg: defaultMiddleware(),
	AuthorizerCfg:    defaultMiddleware(),
	OperationsCfg:    defaultMiddleware(),
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
		return c.Status(500).SendString("Middleware not implemented")
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
	middlewareConf[AuthenticatorCfg] = defaultMiddleware()
	return middlewareConf[AuthenticatorCfg].MiddlewareFunc()
}

// Authorizer returns the Authorizer middleware
func Authorizer() fiber.Handler {
	if mw, ok := middlewareConf[AuthorizerCfg]; ok {
		return mw.MiddlewareFunc()
	}
	middlewareConf[AuthorizerCfg] = defaultMiddleware()
	return middlewareConf[AuthorizerCfg].MiddlewareFunc()
}

// Operations returns the Operations middleware
func Operations() fiber.Handler {
	if mw, ok := middlewareConf[OperationsCfg]; ok {
		return mw.MiddlewareFunc()
	}
	middlewareConf[OperationsCfg] = defaultMiddleware()
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
