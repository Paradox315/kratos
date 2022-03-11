package fiber_validate

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/xhttp/binding"
	"github.com/gofiber/fiber/v2"
)

func NewQueryValidator() *QueryValidator {
	return &QueryValidator{}
}

type QueryValidator struct {
}

func (q *QueryValidator) MiddlewareFunc() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req validator
		_ = binding.BindQuery(c, &req)
		if v, ok := req.(validator); ok {
			if err := v.Validate(); err != nil {
				return errors.BadRequest("VALIDATOR", err.Error())
			}
		}
		return c.Next()
	}
}

func (q *QueryValidator) Name() string {
	return "QueryValidator"
}
