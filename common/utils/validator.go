package utils

import (
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

type ErrorMessageShow struct {
	Message string                 `json:"message"`
	Errors  map[int][]ErrorMessage `json:"errors"`
}

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Validate(c echo.Context, request interface{}) map[int][]ErrorMessage {
	errorMessages := map[int][]ErrorMessage{}

	// Convert the request to a slice of reflect.Values
	requestValues := reflect.ValueOf(request)
	if requestValues.Kind() != reflect.Slice {
		return errorMessages
	}

	for key := 0; key < requestValues.Len(); key++ {
		r := requestValues.Index(key).Interface()

		if err := c.Validate(r); err != nil {
			c.Echo().Logger.Error(err)

			// prepare Error message
			errM := PrepareErrorMessage(err.Error())
			errParts := strings.Split(errM, "=>")

			for i := range errParts {
				if i != 0 && i%2 == 0 {
					breakMessages := strings.SplitN(errParts[i], " ", 2)
					if len(breakMessages) > 0 {
						errorMessages[key] = append(errorMessages[key], ErrorMessage{
							Field:   breakMessages[0],
							Message: breakMessages[1],
						})
					}
				}
			}
		}
	}

	return errorMessages
}

func PrepareErrorMessage(errMessage string) string {
	errM := strings.ReplaceAll(errMessage, " Error:Field validation for ", "=>")
	errM = strings.ReplaceAll(errM, "Key: ", "=>")
	errM = strings.ReplaceAll(errM, "'", "")
	errM = strings.ReplaceAll(errM, "failed on the ", "")
	errM = strings.ReplaceAll(errM, " tag\n", "")
	errM = strings.ReplaceAll(errM, " tag", "")

	return errM
}
