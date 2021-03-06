/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFoundationContext(t *testing.T) {
	c := context.Background()
	c = NewContext(c)
	c = AppendUserIDToContext(c, "ITS ME MARIO")
	c = AppendMethodToContext(c, "grpc", "user.get")
	b := bytes.NewBuffer(nil)

	var log = GetLoggerFromContext(c).
		SetOutput(b).
		WithField("a", "b").WithField("c", "d")
	log.Println("Hello, world")
	assert.Contains(t, b.String(), "Hello, world")
	assert.Contains(t, b.String(), "ITS ME MARIO")
	assert.Contains(t, b.String(), GetRequestIDFromContext(c))
	assert.Contains(t, "ITS ME MARIO", GetUserIDFromContext(c))
}
