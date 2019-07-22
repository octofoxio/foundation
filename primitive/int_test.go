package primitivepb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInt(t *testing.T) {
	data := NewInt(1)
	v, err := data.Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, v)

	var intToScan Int
	err = intToScan.Scan(2)
	assert.NoError(t, err)

}
