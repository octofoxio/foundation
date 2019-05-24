/*
 * Copyright (c) 2019. Octofox.io
 */

package errors

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewGlobErr(t *testing.T) {
	err := NewGlobErr(ErrorTypeInternal, "internal error").WithDetail("TEST").WithDebug("Oh test tes")
	assert.Equal(t, []string{"TEST"}, err.GetDetail())
	assert.Equal(t, "internal error", err.Error())
	assert.Equal(t, ErrorTypeInternal, err.Type())
	assert.Equal(t, "Oh test tes", err.debug)
	assert.Equal(t, "Oh test tes", err.GetDebug())
}
