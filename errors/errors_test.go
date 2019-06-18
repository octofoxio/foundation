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

	b, _ := json.Marshal(err.WithDetail("marshal test"))
	t.Log(string(b))
}
