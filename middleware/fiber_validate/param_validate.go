package fiber_validate

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/xhttp/binding"
	"github.com/gofiber/fiber/v2"
)

func NewParamsValidator() *ParamsValidator {
	return &ParamsValidator{}
}

type ParamsValidator struct {
}

func (p *ParamsValidator) MiddlewareFunc() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req validator
		_ = binding.BindParams(c, req)
		if v, ok := req.(validator); ok {
			if err := v.Validate(); err != nil {
				return errors.BadRequest("VALIDATOR", err.Error())
			}
		}
		return c.Next()
	}
}

func (p *ParamsValidator) Name() string {
	return "ParamValidator"
}
