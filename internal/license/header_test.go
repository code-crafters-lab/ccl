package license

import (
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	newHeader := NewHeader(
		WithVersion(1, 1),
		WithConstantPool(math.MaxInt16),
		WithPayload(math.MaxInt32),
		WithSecurity(math.MaxInt32),
		WithTotal(math.MaxInt32),
	)
	file, err := os.Create("test.lic")
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.Write(newHeader.Bytes())
	assert.Nil(t, err)

	parseHeader, err := ParseHeader([20]byte(newHeader.Bytes()))
	assert.Nil(t, err)
	assert.Equal(t, newHeader, parseHeader)
}
