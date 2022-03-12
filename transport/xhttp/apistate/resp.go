package apistate

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/gofiber/fiber/v2"
)

// Resp	represents the response of the API
type Resp struct {
	Code     int         `json:"code,omitempty"`
	Message  string      `json:"message,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
	Error    interface{} `json:"error,omitempty"`
}

// WithCode sets the response code
func (r *Resp) WithCode(code int) *Resp {
	r.Code = code
	return r
}

// WithMessage sets the message of the response
func (r *Resp) WithMessage(msg string) *Resp {
	r.Message = msg
	return r
}

// WithData sets the data of the response
func (r *Resp) WithData(data interface{}) *Resp {
	r.Metadata = data
	return r
}

// WithError sets the error
func (r *Resp) WithError(err interface{}) *Resp {
	if err, ok := err.(*errors.Error); ok {
		r.Error = err
	} else {
		r.Error = err.Error()
	}
	return r
}

// Send sends the response
func (r *Resp) Send(c *fiber.Ctx) error {
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
func Success() *Resp {
	return &Resp{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
	}
}

// Error response
func Error() *Resp {
	return &Resp{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}
}

// AuthError response
func AuthError() *Resp {
	return &Resp{
		Code:    http.StatusUnauthorized,
		Message: http.StatusText(http.StatusUnauthorized),
	}
}

// InvalidError response
func InvalidError() *Resp {
	return &Resp{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest) + " - Validation Error",
	}
}
