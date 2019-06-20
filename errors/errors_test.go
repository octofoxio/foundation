/*
 * Copyright (c) 2019. Octofox.io
 */

package errors

import (
	"encoding/json"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestNewGlobErr(t *testing.T) {

	err := NewGlobErr(ErrorTypeInternal, "internal error").WithDetail("TEST").WithDebug("Oh test tes")
	assert.EqualValues(t, []string{"TEST"}, err.GetDetail())
	assert.EqualValues(t, "internal error", err.Error())
	assert.EqualValues(t, ErrorTypeInternal, err.Type())
	assert.EqualValues(t, "Oh test tes", err.debug)
	assert.EqualValues(t, "Oh test tes", err.GetDebug())
	err = err.WithDetail("marshal test")
	b, _ := json.Marshal(err)
	t.Log(string(b))

	var merr Error
	_ = json.Unmarshal(b, &merr)
	assert.EqualValues(t, merr.code, err.code)
	assert.EqualValues(t, merr.message, err.message)
	assert.EqualValues(t, merr.detail, err.detail)
}
