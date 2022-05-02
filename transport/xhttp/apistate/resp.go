package apistate

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// Resp	represents the response of the API
type Resp[T any] struct {
	Code     int32  `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	Metadata T      `json:"metadata,omitempty"`
	Error    any    `json:"error,omitempty"`
}

// WithCode sets the response code
func (r *Resp[T]) WithCode(code int32) *Resp[T] {
	r.Code = code
	return r
}

// WithMessage sets the message of the response
func (r *Resp[T]) WithMessage(msg string) *Resp[T] {
	r.Message = msg
	return r
}

// WithData sets the data of the response
func (r *Resp[T]) WithData(data T) *Resp[T] {
	r.Metadata = data
	return r
}

// WithError sets the error
func (r *Resp[T]) WithError(err any) *Resp[T] {
	if err, ok := err.(*errors.Error); ok {
		r.Error = err
		return r
	}
	r.Error = err
	r.Message = err.(error).Error()
	return r
}

// Send sends the response
func (r *Resp[T]) Send(c *fiber.Ctx) error {
	if err, ok := r.Error.(*errors.Error); ok {
		if utils.StatusMessage(int(err.Code)) == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
		return c.Status(int(err.Code)).JSON(err)
	}

	if message := utils.StatusMessage(int(r.Code)); message == "" {
		r.Message = message
		return c.Status(fiber.StatusInternalServerError).JSON(r)
	}
	if r.Message == "" {
		r.Message = utils.StatusMessage(int(r.Code))
	}
	return c.Status(int(r.Code)).JSON(r)
}

// Success response
func Success[T any]() *Resp[T] {
	return &Resp[T]{
		Code:    fiber.StatusOK,
		Message: "success",
	}
}

// Error response
func Error[T any]() *Resp[T] {
	return &Resp[T]{
		Code:    fiber.StatusInternalServerError,
		Message: utils.StatusMessage(fiber.StatusInternalServerError),
	}
}

// InvalidError response
func InvalidError[T any]() *Resp[T] {
	return &Resp[T]{
		Code:    fiber.StatusBadRequest,
		Message: utils.StatusMessage(fiber.StatusBadRequest),
	}
}
