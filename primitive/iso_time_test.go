package primitivepb

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestISOTime_Time(t *testing.T) {
	now := time.Now()
	data := NewISOTime(now)
	assert.EqualValues(t, data.Time().Unix(), now.Unix())
}
