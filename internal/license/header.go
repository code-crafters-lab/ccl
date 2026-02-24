package license

import (
	"encoding/binary"
	"fmt"
)

type Header interface {
	Raw
	// GetMagicNumber 标识文件类型
	GetMagicNumber() [4]byte
	// GetMinorVersion 次要版本号
	GetMinorVersion() uint8
	// GetMajorVersion 主要版本号
	GetMajorVersion() uint8
	// GetReserved 保留字节
	GetReserved() [2]byte
	// GetPayloadLength 授权信息块字节长度
	GetPayloadLength() uint32
	// GetSecurityLength 安全块长度
	GetSecurityLength() uint32
	// GetVendor 获取供应商标识
	GetVendor() [4]byte
	// GetVendorString 获取供应商标识字符串
	GetVendorString() string
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
	// 保留字节（2 字节）
	reserved [2]byte
	// 授权信息块长度（4 字节，32 位无符号整数）
	payloadLength [4]byte
	// 安全块长度（4 字节，32 位无符号整数）
	securityLength [4]byte
	// 4 字节厂商标记，全为 0 表示无厂商
	vendor [4]byte
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

func (h *header) GetReserved() [2]byte {
	return h.reserved
}

func (h *header) GetPayloadLength() uint32 {
	return binary.BigEndian.Uint32(h.payloadLength[:])
}

func (h *header) GetSecurityLength() uint32 {
	return binary.BigEndian.Uint32(h.securityLength[:])
}

func (h *header) GetVendor() [4]byte {
	return h.vendor
}

func (h *header) GetVendorString() string {
	if h.vendor == [4]byte{} {
		return ""
	}
	return fmt.Sprintf("%x", h.vendor)
}

func (h *header) RawBytes() ([]byte, error) {
	result := make([]byte, 20)
	for i, b := range h.magicNumber {
		result[i] = b
	}
	result[4] = h.minorVersion
	result[5] = h.majorVersion
	for i, b := range h.reserved {
		result[6+i] = b
	}
	for i, b := range h.payloadLength {
		result[8+i] = b
	}
	for i, b := range h.securityLength {
		result[12+i] = b
	}
	for i, b := range h.vendor {
		result[16+i] = b
	}
	return result, nil
}

type HeaderOptions func(h *header)

func HeaderWithVersion(major, minor uint8) HeaderOptions {
	return func(h *header) {
		h.minorVersion = minor
		h.majorVersion = major
	}
}

func HeaderWithPayload(length uint32) HeaderOptions {
	return func(h *header) {
		binary.BigEndian.PutUint32(h.payloadLength[:], length)
	}
}

func HeaderWithSecurity(length uint32) HeaderOptions {
	return func(h *header) {
		binary.BigEndian.PutUint32(h.securityLength[:], length)
	}
}

func HeaderWithVendor(vendor [4]byte) HeaderOptions {
	return func(h *header) {
		h.vendor = vendor
	}
}

// HeaderFrom 从一个 20 字节的数组中解析出 Header 信息
// 它会严格按照结构体定义的字段顺序和大小进行解析
func HeaderFrom(bytes [20]byte) (Header, error) {
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

	// Reserved (2 bytes, 6-7)
	copy(h.reserved[:], bytes[6:7])

	// Payload Length (4 bytes, 8-11)
	copy(h.payloadLength[:], bytes[8:12])

	// Security Length (4 bytes, 12-15)
	copy(h.securityLength[:], bytes[12:16])

	// Total Length (4 bytes, 16-19)
	copy(h.vendor[:], bytes[16:20])

	// 4. 返回解析成功的 header 实例
	return h, nil
}
