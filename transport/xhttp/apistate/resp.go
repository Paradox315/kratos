package ApiState

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// RespOption defines the options for the response
type RespOption func(*Resp)

// Resp	represents the response of the API
type Resp struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
	Err  []string    `json:"errors,omitempty"`
}

// SetCode sets the response code
func (r *Resp) SetCode(code int) *Resp {
	r.Code = code
	return r
}

// SetMsg sets the message of the response
func (r *Resp) SetMsg(msg string) *Resp {
	r.Msg = msg
	return r
}

// SetData sets the data of the response
func (r *Resp) SetData(data interface{}) *Resp {
	r.Data = data
	return r
}

// SetError sets the error
func (r *Resp) SetError(err ...string) *Resp {
	r.Err = err
	return r
}

// Send sends the response
func (r Resp) Send(c *fiber.Ctx) error {
	if http.StatusText(r.Code) == "" {
		r.Msg = http.StatusText(http.StatusInternalServerError)
		return c.Status(http.StatusInternalServerError).JSON(r)
	}
	if r.Msg == "" {
		r.Msg = http.StatusText(r.Code)
	}
	return c.Status(r.Code).JSON(r)
}

// SendMessage sends the response with the message
func (r Resp) SendMessage(c *fiber.Ctx, msg string) error {
	r.Msg = msg
	if http.StatusText(r.Code) == "" {
		return c.Status(http.StatusInternalServerError).JSON(r)
	}
	return c.Status(r.Code).JSON(r)
}

// SendData sends the response with the data
func (r Resp) SendData(c *fiber.Ctx, data interface{}) error {
	r.Data = data
	if http.StatusText(r.Code) == "" {
		return c.Status(http.StatusInternalServerError).JSON(r)
	}
	if r.Msg == "" {
		r.Msg = http.StatusText(r.Code)
	}
	return c.Status(r.Code).JSON(r)
}

// SendError sends the response with the error
func (r Resp) SendError(c *fiber.Ctx, err ...string) error {
	r.Err = err
	if http.StatusText(r.Code) == "" {
		return c.Status(http.StatusInternalServerError).JSON(r)
	}
	if r.Msg == "" {
		r.Msg = http.StatusText(r.Code)
	}
	return c.Status(r.Code).JSON(r)
}

// WithCode function sets the response code
func WithCode(code int) RespOption {
	return func(r *Resp) {
		r.Code = code
	}
}

// WithMsg function sets the response message
func WithMsg(msg string) RespOption {
	return func(r *Resp) {
		r.Msg = msg
	}
}

// WithData function sets the response data
func WithData(data interface{}) RespOption {
	return func(r *Resp) {
		r.Data = data
	}
}

// WithErr function sets the response error
func WithErr(err ...string) RespOption {
	return func(r *Resp) {
		r.Err = err
	}
}

// New response
func New(opts ...RespOption) *Resp {
	r := &Resp{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Success response
func Success(opts ...RespOption) *Resp {
	r := &Resp{
		Code: http.StatusOK,
		Msg:  http.StatusText(http.StatusOK),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Error response
func Error(opts ...RespOption) *Resp {
	r := &Resp{
		Code: http.StatusInternalServerError,
		Msg:  http.StatusText(http.StatusInternalServerError),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// AuthError response
func AuthError(opts ...RespOption) *Resp {
	r := &Resp{
		Code: http.StatusUnauthorized,
		Msg:  http.StatusText(http.StatusUnauthorized),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// InvalidError response
func InvalidError(opts ...RespOption) *Resp {
	r := &Resp{
		Code: http.StatusBadRequest,
		Msg:  http.StatusText(http.StatusBadRequest) + " - Validation Error",
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
