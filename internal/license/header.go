package license

import (
	"encoding/binary"
	"fmt"
)

type Header interface {
	// GetMagicNumber 标识文件类型
	GetMagicNumber() [4]byte
	// GetMinorVersion 次要版本号
	GetMinorVersion() uint8
	// GetMajorVersion 主要版本号
	GetMajorVersion() uint8
	// GetConstantPoolCount 常量池计数
	GetConstantPoolCount() uint16
	// GetPayloadOffset 授权信息块偏移
	GetPayloadOffset() uint32
	// GetSecurityOffset 安全块偏移
	GetSecurityOffset() uint32
	// GetTotalLength 文件总长度
	GetTotalLength() uint32

	Bytes() []byte
}

func NewHeader(opts ...HeaderOptions) Header {
	h := &header{
		magicNumber:  [4]byte{0xCC, 0x4C, 0x9A, 0x9B},
		minorVersion: 1,
		majorVersion: 0,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

type header struct {
	// 标识文件类型（4 字节），固定为 0xCC4C9A9B
	magicNumber [4]byte
	// 次要版本号（1 字节，8 位无符号整数）
	minorVersion byte
	// 主要版本号（1 字节，8 位无符号整数）
	majorVersion byte
	// 常量池计数（2 字节，16 位无符号整数）
	constantPoolCount [2]byte
	// 授权信息块偏移（4 字节，32 位无符号整数）
	payloadOffset [4]byte
	// 安全块偏移（4 字节，32 位无符号整数）
	securityOffset [4]byte
	// 文件总长度（4 字节，32 位无符号整数）
	totalLength [4]byte
}

func (h *header) GetMagicNumber() [4]byte {
	return h.magicNumber
}

func (h *header) GetMinorVersion() uint8 {
	return h.minorVersion
}

func (h *header) GetMajorVersion() uint8 {
	return h.majorVersion
}

func (h *header) GetConstantPoolCount() uint16 {
	return binary.BigEndian.Uint16(h.constantPoolCount[:])
}

func (h *header) GetPayloadOffset() uint32 {
	return binary.BigEndian.Uint32(h.payloadOffset[:])
}

func (h *header) GetSecurityOffset() uint32 {
	return binary.BigEndian.Uint32(h.securityOffset[:])
}

func (h *header) GetTotalLength() uint32 {
	return binary.BigEndian.Uint32(h.totalLength[:])
}

func (h *header) Bytes() []byte {
	result := make([]byte, 20)
	for i, b := range h.magicNumber {
		result[i] = b
	}
	result[4] = h.minorVersion
	result[5] = h.majorVersion
	for i, b := range h.constantPoolCount {
		result[6+i] = b
	}
	for i, b := range h.payloadOffset {
		result[8+i] = b
	}
	for i, b := range h.securityOffset {
		result[12+i] = b
	}
	for i, b := range h.totalLength {
		result[16+i] = b
	}
	return result
}

type HeaderOptions func(h *header)

func WithVersion(major, minor uint8) HeaderOptions {
	return func(h *header) {
		h.minorVersion = minor
		h.majorVersion = major
	}
}

func WithConstantPool(count uint16) HeaderOptions {
	return func(h *header) {
		binary.BigEndian.PutUint16(h.constantPoolCount[:], count)
	}
}

func WithPayload(offset uint32) HeaderOptions {
	return func(h *header) {
		binary.BigEndian.PutUint32(h.payloadOffset[:], offset)
	}
}

func WithSecurity(offset uint32) HeaderOptions {
	return func(h *header) {
		binary.BigEndian.PutUint32(h.securityOffset[:], offset)
	}
}

func WithTotal(length uint32) HeaderOptions {
	return func(h *header) {
		binary.BigEndian.PutUint32(h.totalLength[:], length)
	}
}

// ParseHeader 从一个 20 字节的数组中解析出 Header 信息
// 它会严格按照结构体定义的字段顺序和大小进行解析
func ParseHeader(bytes [20]byte) (Header, error) {
	// 1. 创建一个新的 header 实例作为解析目标
	h := &header{}

	// 2. 按顺序解析各个字段
	// Magic Number (4 bytes, 0-3)
	copy(h.magicNumber[:], bytes[0:4])

	// 3. 验证 Magic Number 是否正确
	expectedMagic := [4]byte{0xCC, 0x4C, 0x9A, 0x9B}
	if h.magicNumber != expectedMagic {
		return nil, fmt.Errorf("invalid license header")
	}

	// Minor Version (1 byte, 4)
	h.minorVersion = bytes[4]

	// Major Version (1 byte, 5)
	h.majorVersion = bytes[5]

	// Constant Pool Count (2 bytes, 6-7)
	// 使用 binary.BigEndian.Uint16 来解析这两个字节
	h.constantPoolCount[0] = bytes[6]
	h.constantPoolCount[1] = bytes[7]

	// Payload Offset (4 bytes, 8-11)
	copy(h.payloadOffset[:], bytes[8:12])

	// Security Offset (4 bytes, 12-15)
	copy(h.securityOffset[:], bytes[12:16])

	// Total Length (4 bytes, 16-19)
	copy(h.totalLength[:], bytes[16:20])

	// 4. 返回解析成功的 header 实例
	return h, nil
}
