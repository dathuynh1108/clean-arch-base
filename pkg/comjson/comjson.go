package comjson

import (
	"io"

	json "github.com/bytedance/sonic"
)

var (
	Marshal         func(val interface{}) ([]byte, error)   = json.Marshal
	Unmarshal       func(buf []byte, val interface{}) error = json.Unmarshal
	UnmarshalString func(buf string, val interface{}) error = json.UnmarshalString
	MarshalString   func(val interface{}) (string, error)   = json.MarshalString
	NewEncoder      func(writer io.Writer) json.Encoder     = json.ConfigDefault.NewEncoder
	NewDecoder      func(reader io.Reader) json.Decoder     = json.ConfigDefault.NewDecoder
)
