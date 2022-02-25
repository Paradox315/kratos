package ApiState

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

const (
	SuccessState = "success"
	ErrorState   = "error"
)

// RespOption defines the options for the response
type RespOption func(*Resp)

// Resp	represents the response of the API
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
	Err  []error     `json:"errors"`
}

func (r Resp) Send(c *fiber.Ctx) error {
	return c.JSON(r)
}

func (r Resp) SendMessage(c *fiber.Ctx, msg string) error {
	r.Msg = msg
	return c.JSON(r)
}

func (r Resp) SendData(c *fiber.Ctx, data interface{}) error {
	r.Data = data
	return c.JSON(r)
}

func (r Resp) SendError(c *fiber.Ctx, err ...error) error {
	r.Err = err
	return c.JSON(r)
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
func WithErr(err ...error) RespOption {
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
		Msg:  SuccessState,
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
		Msg:  ErrorState,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
