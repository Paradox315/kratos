package binding

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/form"
	"github.com/go-kratos/kratos/v2/transport/xhttp"
	"github.com/gofiber/fiber/v2"
	"net/url"
)

// BindQuery bind query parameters to target.
func BindQuery(ctx *fiber.Ctx, target interface{}) (err error) {
	if err = ctx.QueryParser(target); err != nil {
		return err
	}
	if err = xhttp.Validate(target); err != nil {
		return err
	}
	return
}

// BindBody bind body parameters to target.
func BindBody(ctx *fiber.Ctx, target interface{}) (err error) {
	if err = ctx.BodyParser(target); err != nil {
		return err
	}
	if err = xhttp.Validate(target); err != nil {
		return err
	}
	return
}

// BindParams bind body parameters to target.
func BindParams(ctx *fiber.Ctx, target interface{}) (err error) {
	vars := make(url.Values, len(ctx.Route().Params))
	for _, k := range ctx.Route().Params {
		vars[k] = []string{ctx.Params(k)}
	}
	if err = encoding.GetCodec(form.Name).Unmarshal([]byte(vars.Encode()), target); err != nil {
		return err
	}
	if err = xhttp.Validate(target); err != nil {
		return err
	}
	return
}
