package license

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	newHeader := NewHeader(
		HeaderWithVersion(1, 1),
		HeaderWithPayload(234),
		HeaderWithSecurity(356),
		HeaderWithVendor([4]byte{0x0, 0x0, 0x0, 0x1}),
	)
	file, err := os.Create("test.lic")
	defer file.Close()
	assert.Nil(t, err)
	bytes, err := newHeader.RawBytes()
	assert.Nil(t, err)
	_, err = file.Write([]byte(base64.StdEncoding.EncodeToString(bytes)))
	assert.Nil(t, err)

	parseHeader, err := HeaderFrom([20]byte(bytes))
	assert.Nil(t, err)
	assert.Equal(t, newHeader, parseHeader)
}
