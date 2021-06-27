package storage

import (
	"encoding/json"
	"io/ioutil"
)

func decode(filename string, output interface{}) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, output); err != nil {
		return err
	}

	return nil
}
