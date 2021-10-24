package http

import (
	"encoding/json"
	"net/http"

	"github.com/why444216978/go-util/validate"
)

func ParseAndValidateBody(req *http.Request, target interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(target); err != nil {
		return err
	}
	if err := validate.Validate(target); err != nil {
		return err
	}

	return nil
}
