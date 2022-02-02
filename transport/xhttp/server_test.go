package xhttp

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	httpSrv := NewServer(
		Address(":19000"),
		Router(func(r fiber.Router) {
			r.Get("api", func(ctx *fiber.Ctx) error {
				return ctx.JSON(fiber.Map{"status": http.StatusOK, "path": ctx.Path()})
			})
		}),
		Middlewares(
			recover.New(recover.Config{
				EnableStackTrace: true,
			}),
			logger.New(logger.ConfigDefault),
		),
	)
	httpSrv.WithRouter(func(r fiber.Router) {
		r.Get("hello", func(ctx *fiber.Ctx) error {
			return ctx.SendString("hello")
		})
		r.Get("hello/:name", func(ctx *fiber.Ctx) error {
			return ctx.SendString("hello" + ctx.Params("name"))
		})
	})
	app := kratos.New(
		kratos.Name("fiber"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
