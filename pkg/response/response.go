package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Message struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

var (
	InternalServerError = Message{
		Message: "Internal Server Error",
		Details: nil,
	}
	InputErrResponse = Message{
		Message: "Input error",
		Details: nil,
	}
)

func InputFieldError(err error) Message {
	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		return InputErrResponse
	}

	var errs []string
	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, err.Field()+" is invalid ("+err.Tag()+")")
	}

	return Message{
		Message: "Input error",
		Details: errs,
	}
}
