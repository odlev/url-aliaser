// Package api is a nice package
package api

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"` //Error, Ok
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK = "OK" 
	StatusError = "Error" 
)

func OK() Response {
	return Response {
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error: msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMessage []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required": 
			errMessage = append(errMessage, fmt.Sprintf("field %s is a required field, you are miss it", err.Field()))
		case "url":
			errMessage = append(errMessage, fmt.Sprintf("field %s is not a valid url, try again", err.Field()))
		default:
			errMessage = append(errMessage, fmt.Sprintf("field %s in not a valiv", err.Field()))
		}
		
	}
	return Response{
		Status: StatusError,
		Error: strings.Join(errMessage, ", "),
	}
}
