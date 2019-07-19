/*
 * Copyright (c) 2019. Octofox.io
 */

package primitivev1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	data := NewString("Hello")
	v, err := data.Value()
	assert.NoError(t, err)
	assert.EqualValues(t, v, "Hello")
	assert.False(t, data.GetIsNull())

	var stringToScan String
	err = stringToScan.Scan("World")
	assert.NoError(t, err)
	assert.EqualValues(t, stringToScan.GetV(), "World")

}
