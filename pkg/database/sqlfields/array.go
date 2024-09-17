package sqlfields

import (
	"database/sql/driver"

	"github.com/dathuynh1108/clean-arch-base/pkg/comjson"
)

type Array[T any] []T

func (array Array[T]) Value() (driver.Value, error) {
	bytes, err := comjson.Marshal(array)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (array *Array[T]) Scan(input interface{}) error {
	inputBytes := input.([]byte)
	if len(inputBytes) == 0 {
		inputBytes = []byte("[]")
	}
	return comjson.Unmarshal(inputBytes, array)
}
