package codec

import (
	json "github.com/json-iterator/go"
)

type Codec interface {
	Encode(data interface{}) ([]byte, error)
	Decode(data []byte, i interface{}) error
}

type JSONCodec struct{}

var _ Codec = (*JSONCodec)(nil)

func (c JSONCodec) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (c JSONCodec) Decode(data []byte, dst interface{}) error {
	return json.Unmarshal(data, dst)
}
