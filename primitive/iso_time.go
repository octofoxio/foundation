/*
 * Copyright (c) 2019. Octofox.io
 */

package primitivev1

import (
	"database/sql/driver"
	"github.com/octofoxio/foundation/errors"
	"reflect"
	"time"
)

func NewISOTime(t time.Time) *ISOTime {
	return &ISOTime{
		V: t.Format(time.RFC3339Nano),
	}
}

func (i *ISOTime) Time() time.Time {
	t, _ := time.Parse(time.RFC3339Nano, i.V)
	return t
}

func (m ISOTime) Value() (driver.Value, error) {
	if m.GetIsNull() {
		return nil, nil
	}
	return m.GetV(), nil
}

func (m *ISOTime) Scan(src interface{}) error {
	if src == nil {
		m.IsNull = true
		return nil
	}

	k := reflect.TypeOf(src).Kind()
	if k == reflect.String {
		_, err := time.Parse(time.RFC3339Nano, src.(string))
		if err != nil {
			return err
		}
		m.V = src.(string)
		return nil
	}

	return errors.New(errors.ErrorTypeInternal, "int failed to scan")
}
