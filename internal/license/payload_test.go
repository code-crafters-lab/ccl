package license

import (
	"testing"
	"time"

	"github.com/code-crafters-lab/ccl/pkg/grpc/license"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPayload_Empty(t *testing.T) {
	pl := NewPayload()
	bytes, err := pl.RawBytes()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(bytes))
}

func pt[T any](v T) *T {
	return &v
}

func TestPayload(t *testing.T) {
	now := time.Now()
	issuedAt := timestamppb.New(now)
	expiredAt := timestamppb.New(now.Add(24 * time.Hour))

	securityBlock := &license.SecurityBlock{}

	var (
		a1         = &anypb.Any{}
		extensions = []*anypb.Any{a1}
	)

	_ = anypb.MarshalFrom(a1, securityBlock, proto.MarshalOptions{})

	t.Log(a1.GetTypeUrl())
	t.Log(a1.MessageName())
	t.Log(a1.MessageIs(&license.SecurityBlock{}))
	t.Log(a1.MessageIs(&license.HardwareLicenseSpec{}))

	p := &license.Payload{
		LicenseId:   pt("1"),
		ProductId:   pt("CC"),
		LicenseMode: pt(license.Mode_MODE_SOFTWARE),
		LicenseType: pt(license.Type_TYPE_INDIVIDUAL),
		IssuedAt:    issuedAt,
		ExpiredAt:   expiredAt,
		Functions: &license.Payload_FunctionsValue{
			FunctionsValue: "15",
		},
		Extensions: extensions,
		Spec: &license.Payload_SoftwareSpec{
			SoftwareSpec: &license.SoftwareLicenseSpec{
				DeviceFingerprint: []string{"45555"},
			},
		},
	}
	pl := NewPayload(
		WithData(p),
		WithLicense(nil),
		WithDeadline("2025-12-31"),
		WithVersion("1.1"),
	)
	bytes, err := pl.RawBytes()
	assert.NoError(t, err)
	t.Logf("\n%s\n", string(bytes))
}

func TestPayload_Header(t *testing.T) {
	h := NewHeader(
		HeaderWithVersion(1, 2),
	)
	pl := NewPayload(WithData(h))
	bytes, err := pl.RawBytes()
	assert.NoError(t, err)
	assert.Equal(t, 20, len(bytes))
	assert.Equal(t, []byte{0xCC, 0x4C, 0x9A, 0x9B}, bytes[0:4])
	assert.Equal(t, uint8(2), bytes[4])
	assert.Equal(t, uint8(1), bytes[5])
}
