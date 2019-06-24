/*
 * Copyright (c) 2019. Octofox.io
 */

package fs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewLocalFileStorage(t *testing.T) {
	wd, _ := os.Getwd()
	local := NewLocalFileStorage(wd)

	u, err := local.GetObjectURL("./storage.go")
	assert.NoError(t, err)
	assert.NotEmpty(t, u)
	t.Log(u)

}
