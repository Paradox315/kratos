package apistate

import (
	"net/http"
	"reflect"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/gofiber/fiber/v2"
)

// Resp	represents the response of the API
type Resp[T any] struct {
	Code     int         `json:"code,omitempty"`
	Message  string      `json:"message,omitempty"`
	Metadata T           `json:"metadata,omitempty"`
	Error    interface{} `json:"error,omitempty"`
}

// WithCode sets the response code
func (r *Resp[T]) WithCode(code int) *Resp[T] {
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
	errMsg := reflect.ValueOf(err).MethodByName("Error").Call(nil)
	r.Error = errMsg[0].Interface()
	return r
}

// Send sends the response
func (r *Resp[T]) Send(c *fiber.Ctx) error {
	if err, ok := r.Error.(*errors.Error); ok {
		if http.StatusText(int(err.Code)) == "" {
			return c.Status(http.StatusInternalServerError).JSON(err)
		}
		return c.Status(int(err.Code)).JSON(err)
	}

	if http.StatusText(r.Code) == "" {
		r.Message = http.StatusText(http.StatusInternalServerError)
		return c.Status(http.StatusInternalServerError).JSON(r)
	}
	if r.Message == "" {
		r.Message = http.StatusText(r.Code)
	}
	return c.Status(r.Code).JSON(r)
}

// Success response
func Success[T any]() *Resp[T] {
	return &Resp[T]{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
	}
}

// Error response
func Error[T any]() *Resp[T] {
	return &Resp[T]{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}
}
