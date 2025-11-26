package license

import (
	"fmt"
	"os"
	"testing"

	"github.com/code-crafters-lab/ccl/pkg/grpc/license"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestCreateLicense(t *testing.T) {
	hexStr := "CC4C9A9B"
	var magicNumber uint32

	// 使用 %x 格式符来扫描十六进制整数
	// &uint32Value 是接收扫描结果的变量地址
	_, err := fmt.Sscanf(hexStr, "%x", &magicNumber)

	// 1. 创建LicenseFile
	//header := &license.Header{
	//	MagicNumber: &magicNumber,
	//}
	payload := &license.Payload{}
	//licenseFile := &license.LicenseFile{
	//	Header: header,
	//}

	// 2. 序列化（proto.Marshal）
	payloadBytes, err := proto.Marshal(payload)
	assert.Nil(t, err)

	file, err := os.Create("test.lic")
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.Write(payloadBytes)
	assert.Nil(t, err)
}
