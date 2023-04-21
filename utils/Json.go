package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kando-backend/httpErrors"
	"net/http"
	"strings"
)

func WriteJson(w http.ResponseWriter, src interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	err := enc.Encode(src)

	if err != nil {
		return httpErrors.InternalServerError().WithMessage("Error encoding JSON")
	}

	w.Header().Set("Content-Type", "application/json")

	return nil
}

func ReadJson(r io.Reader, dst interface{}) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return httpErrors.BadRequest().WithMessage(msg)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return httpErrors.BadRequest().WithMessage(msg)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return httpErrors.BadRequest().WithMessage(msg)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return httpErrors.BadRequest().WithMessage(msg)

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return httpErrors.BadRequest().WithMessage(msg)

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return httpErrors.PayloadTooLarge().WithMessage(msg)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return httpErrors.BadRequest().WithMessage(msg)
	}

	return nil
}
