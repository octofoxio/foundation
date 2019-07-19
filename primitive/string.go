package primitivev1

import (
	"database/sql/driver"
	"github.com/octofoxio/foundation/errors"
)

func NewString(v string) *String {
	return &String{
		IsNull: false,
		V:      v,
	}
}

func NewNullableString(v string) *String {
	return &String{
		IsNull: len(v) == 0,
		V:      v,
	}
}

func (m String) Value() (driver.Value, error) {
	if m.GetIsNull() {
		return nil, nil
	}
	return m.GetV(), nil
}

func (m *String) Scan(src interface{}) error {
	if src == nil {
		m.IsNull = true
		return nil
	}
	if s, ok := src.(string); ok {
		m.V = s
		return nil
	}
	return errors.New(errors.ErrorTypeInternal, "int failed to scan")
}
