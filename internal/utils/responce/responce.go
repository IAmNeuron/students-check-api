package responce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type Respoce struct {
	Status string
	Error  string
}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("conteny-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)

}

func GeneralError(err error) Respoce {
	return Respoce{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidattionError(errs validator.ValidationErrors) Respoce {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Respoce{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}

}
