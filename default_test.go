package foundation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvString(t *testing.T) {
	var x = EnvString("K", "x")
	assert.Equal(t, x, "x")
}
