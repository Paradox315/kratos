package ApiState

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
	"testing"
)

func TestResp_Send(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return Success().SetData(map[string]interface{}{
			"name": "John",
		}).SetMsg("Hello, John!").SetError(errors.New("valid error").Error()).Send(c)
	})

	log.Fatal(app.Listen(":3000"))
}
