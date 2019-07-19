package primitivev1

import (
	"database/sql/driver"
	"github.com/octofoxio/foundation/errors"
	"reflect"
)

func NewInt(i int64) *Int {
	return &Int{
		V: i,
	}
}

func NewNullableInt(i int64) *Int {
	return &Int{
		V:      i,
		IsNull: i == 0,
	}
}

func (m Int) Value() (driver.Value, error) {
	if m.GetIsNull() {
		return nil, nil
	}
	return m.GetV(), nil
}

func (m *Int) Scan(src interface{}) error {
	if src == nil {
		m.IsNull = true
		return nil
	}

	k := reflect.TypeOf(src).Kind()
	if k == reflect.Int || k == reflect.Int64 || k == reflect.Int32 {
		m.V = reflect.ValueOf(src).Int()
		return nil
	}

	return errors.New(errors.ErrorTypeInternal, "int failed to scan")
}
