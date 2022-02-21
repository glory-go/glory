package utils

import jsoniter "github.com/json-iterator/go"

func ConvertInto(from, to interface{}) error {
	data, err := jsoniter.Marshal(from)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(data, to)
}
