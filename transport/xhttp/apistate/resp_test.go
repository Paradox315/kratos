package apistate

import (
	"errors"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/gofiber/fiber/v2"
	"strings"

	"log"
	"testing"
)

type T struct {
	errs []string
}

func (t T) Error() string {
	return strings.Join(t.errs, ";")
}

func TestResp_Send(t *testing.T) {
	app := fiber.New()

	app.Get("/err1", func(c *fiber.Ctx) error {
		return Error[any]().WithError(kerrors.BadRequest("test", "test error")).Send(c)
	})

	app.Get("/err1/1", func(c *fiber.Ctx) error {
		return Error[any]().WithError(kerrors.New(20000, "test", "test error")).Send(c)
	})

	app.Get("/err2", func(c *fiber.Ctx) error {
		return Error[any]().WithError(errors.New("hello noonoo")).Send(c)
	})

	app.Get("/success", func(c *fiber.Ctx) error {
		return Success[fiber.Map]().WithData(fiber.Map{
			"name": "kratos success",
		}).Send(c)
	})

	log.Fatal(app.Listen(":3000"))
}
