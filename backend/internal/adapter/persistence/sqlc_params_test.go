package persistence

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt32Param(t *testing.T) {
	assert.Equal(t, int32(0), int32Param(0))
	assert.Equal(t, int32(20), int32Param(20))
	assert.Equal(t, int32(0), int32Param(-1), "負値は0にクランプする")
	assert.Equal(t, int32(math.MaxInt32), int32Param(math.MaxInt32+1),
		"int32の範囲を超える値はサイレントにオーバーフローさせずint32の最大値にクランプする")
}
