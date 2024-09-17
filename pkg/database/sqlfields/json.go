package sqlfields

import (
	"database/sql/driver"

	"github.com/dathuynh1108/clean-arch-base/pkg/comjson"
)

type JSONAdapter[T any] struct {
	Object T
}

func (a JSONAdapter[T]) Value() (driver.Value, error) {
	return comjson.Marshal(a.Object)

}

func (a *JSONAdapter[T]) Scan(input interface{}) error {
	if input == nil {
		return nil
	}
	return comjson.Unmarshal(input.([]byte), &a.Object)
}
