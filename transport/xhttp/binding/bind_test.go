package binding

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"testing"
)

type HelloRequest struct {
	Name  string `json:"name" form:"name" validate:"required"`
	ID    string `json:"id" form:"id" validate:"required"`
	Phone string `json:"phone" form:"phone"`
}

func TestBindQuery(t *testing.T) {
	app := fiber.New()
	app.Get("hello", func(ctx *fiber.Ctx) error {
		var in HelloRequest
		if err := BindQuery(ctx, &in); err != nil {
			return err
		}
		return ctx.JSON(in)
	})
	log.Fatal(app.Listen(":19000"))
}

func TestBindBody(t *testing.T) {
	app := fiber.New()
	app.Get("hello", func(ctx *fiber.Ctx) error {
		var in HelloRequest
		if err := BindBody(ctx, &in); err != nil {
			return err
		}
		return ctx.JSON(in)
	})
	log.Fatal(app.Listen(":19000"))
}

func TestBindParams(t *testing.T) {
	app := fiber.New()
	app.Get("hello/:name/:phone/:id", func(ctx *fiber.Ctx) error {
		var in HelloRequest
		if err := BindParams(ctx, &in); err != nil {
			return err
		}
		return ctx.JSON(in)
	})
	log.Fatal(app.Listen(":19000"))
}
