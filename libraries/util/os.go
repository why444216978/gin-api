package util

import (
	"encoding/json"
	"io/ioutil"
)

func UnmarshalJsonFromFile(path string, v interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}
