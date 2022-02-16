package binding

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/form"
	"github.com/gofiber/fiber/v2"
	"net/url"
)

// BindQuery bind query parameters to target.
func BindQuery(ctx *fiber.Ctx, target interface{}) (err error) {
	return ctx.QueryParser(target)
}

// BindBody bind body parameters to target.
func BindBody(ctx *fiber.Ctx, target interface{}) (err error) {
	return ctx.BodyParser(target)
}

// BindParams bind body parameters to target.
func BindParams(ctx *fiber.Ctx, target interface{}) (err error) {
	vars := make(url.Values, len(ctx.Route().Params))
	for _, k := range ctx.Route().Params {
		vars[k] = []string{ctx.Params(k)}
	}
	return encoding.GetCodec(form.Name).Unmarshal([]byte(vars.Encode()), target)
}
