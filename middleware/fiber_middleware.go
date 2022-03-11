package middleware

import (
	validate "github.com/go-kratos/kratos/v2/middleware/fiber_validate"
	"github.com/gofiber/fiber/v2"
)

// Constants for the middleware
const (
	SupportPackageIsVersion1 = true
	AuthenticatorCfg         = "Authenticator"
	AuthorizerCfg            = "Authorizer"
	OperationsCfg            = "Operations"
	ValidatorCfg             = "Validator"
)
const (
	BodyTyp = iota + 1
	QueryTyp
	ParamTyp
)

var middlewareConf = map[string]FiberMiddleware{
	AuthenticatorCfg: defaultMiddleware(),
	AuthorizerCfg:    defaultMiddleware(),
	OperationsCfg:    defaultMiddleware(),
	ValidatorCfg:     validate.NewBodyValidator(),
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
func RegisterMiddleware(name string, mw FiberMiddleware) {
	if middlewareConf == nil {
		middlewareConf = make(map[string]FiberMiddleware)
	}
	middlewareConf[name] = mw
}

// Authenticator returns the Authenticator middleware
func Authenticator() fiber.Handler {
	if mw, ok := middlewareConf[AuthenticatorCfg]; ok {
		return mw.MiddlewareFunc()
	}
	return defaultMiddleware().MiddlewareFunc()
}

// Authorizer returns the Authorizer middleware
func Authorizer() fiber.Handler {
	if mw, ok := middlewareConf[AuthorizerCfg]; ok {
		return mw.MiddlewareFunc()
	}
	return defaultMiddleware().MiddlewareFunc()
}

// Operations returns the Operations middleware
func Operations() fiber.Handler {
	if mw, ok := middlewareConf[OperationsCfg]; ok {
		return mw.MiddlewareFunc()
	}
	return defaultMiddleware().MiddlewareFunc()
}

// Validator returns the Validator middleware
func Validator(typ ...int) fiber.Handler {
	if len(typ) > 0 {
		switch typ[0] {
		case QueryTyp:
			return validate.NewQueryValidator().MiddlewareFunc()
		case BodyTyp:
			return validate.NewBodyValidator().MiddlewareFunc()
		case ParamTyp:
			return validate.NewParamsValidator().MiddlewareFunc()
		default:
			return defaultMiddleware().MiddlewareFunc()
		}
	}
	if mw, ok := middlewareConf[ValidatorCfg]; ok {
		return mw.MiddlewareFunc()
	}
	return defaultMiddleware().MiddlewareFunc()
}

// CustomMiddleware returns a custom middleware with your config key
func CustomMiddleware(name string) fiber.Handler {
	if mw, ok := middlewareConf[name]; ok {
		return mw.MiddlewareFunc()
	}
	return defaultMiddleware().MiddlewareFunc()
}
