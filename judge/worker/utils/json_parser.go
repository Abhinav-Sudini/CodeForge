package utils

import (
	"encoding/json"
	"errors"
	"io"
)

func CustomJsonUnMarshal(stream io.ReadCloser,v any) error {
	dec := json.NewDecoder(stream)
	dec.DisallowUnknownFields()

	err:= dec.Decode(v)
	if err != nil {
		return err
	}

	// optional extra check
	if dec.More() {
			return errors.New("[json unmarshal] extra data found in json")
	}
	return nil
}
