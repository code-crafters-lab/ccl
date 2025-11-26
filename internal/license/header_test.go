package license

import (
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	newHeader := NewHeader(
		HeaderWithVersion(1, 1),
		HeaderWithConstantPool(math.MaxInt16),
		HeaderWithPayload(math.MaxInt32),
		HeaderWithSecurity(math.MaxInt32),
		HeaderWithTotal(math.MaxInt32),
	)
	file, err := os.Create("test.lic")
	defer file.Close()
	assert.Nil(t, err)
	bytes, err := newHeader.RawBytes()
	assert.Nil(t, err)
	_, err = file.Write(bytes)
	assert.Nil(t, err)

	parseHeader, err := HeaderFrom([20]byte(bytes))
	assert.Nil(t, err)
	assert.Equal(t, newHeader, parseHeader)
}
