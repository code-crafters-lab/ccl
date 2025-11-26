package license

import (
	"bytes"
	"encoding/json"
	"encoding/pem"
	"reflect"

	"github.com/code-crafters-lab/ccl/internal/license/transform"
	"google.golang.org/protobuf/proto"
)

type Payload interface {
	Raw
	// ID 许可证唯一标识
	//ID() string
	//// IssuedAt 颁发时间
	//IssuedAt() time.Time
	//// ExpiredAt 过期时间
	//ExpiredAt() time.Time
	// Type Header() Header
	//// IsValid 2. 时效性校验
	//IsValid() bool   // 当前是否有效（未过期且未吊销）
	//IsExpired() bool // 是否已过期
	//
	//// AllowedFeatures 3. 功能限制
	//AllowedFeatures() []string            // 获取所有允许的功能
	//IsFeatureAllowed(feature string) bool // 检查某个功能是否允许
	//
	//// BoundDevices 4. 设备绑定
	//BoundDevices() []string             // 获取绑定的设备ID列表
	//IsDeviceBound(deviceID string) bool // 检查设备是否已绑定
	//BindDevice(deviceID string) error   // 绑定新设备（可能有数量限制）
	//
	//// Signature 5. 安全相关
	//Signature() string                      // 获取许可证签名
	//VerifySignature(publicKey string) error // 验证签名有效性
	//
	//// Revoke 6. 状态管理
	//Revoke() error   // 吊销许可证
	//IsRevoked() bool // 检查是否已吊销
}

type CustomPEMHeader func(headers map[string]string)

const (
	HeaderProcType           = "Proc-Type"
	HeaderVersion            = "Version"
	HeaderDeadline           = "Deadline"
	HeaderLicenseRestriction = "License-Restriction"
)

func innerHeader(info *payload[any]) CustomPEMHeader {
	return func(headers map[string]string) {
		//if info.fileName != "" {
		//	headers[HeaderProcType] = info.fileName + " Licence"
		//}

		if info.version != "" {
			headers[HeaderVersion] = info.version
		}

		if info.deadline != "" {
			headers[HeaderDeadline] = info.deadline
		}

		if info.restriction != "" {
			headers[HeaderLicenseRestriction] = info.restriction
		}
	}
}

type payload[T interface{}] struct {
	data T

	pemType       string // pem 类型
	version       string // 授权版本
	deadline      string // 授权截止时间戳
	restriction   string // 限制使用说明
	customHeaders []CustomPEMHeader
	trans         []transform.Transform
}

type PayloadOption func(p *payload[any])

func NewPayload(opts ...PayloadOption) Payload {
	p := &payload[any]{
		customHeaders: make([]CustomPEMHeader, 0),
		trans:         make([]transform.Transform, 0),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (pl *payload[any]) RawBytes() ([]byte, error) {
	var (
		x      interface{} = pl.data
		result []byte
		err    error
	)
	if pl == nil || x == nil {
		result = []byte{}
	}
	// 使用类型断言来判断 data 的具体类型
	switch v := x.(type) {
	case string:
		result = []byte(v)
	case []byte:
		result = v
	case Raw:
		result, err = v.RawBytes()
	// 处理 Protobuf 消息
	case proto.Message:
		result, err = proto.Marshal(v)
	// 对于其他类型，尝试 JSON 序列化
	default:
		if v == nil {
			result = []byte{}
		} else {
			result, err = json.Marshal(v)
		}
	}

	if err != nil {
		return nil, err
	}

	// 数据转换
	for _, tran := range pl.trans {
		result, err = tran(result)
		if err != nil {
			return nil, err
		}
	}
	// pem 自动输出
	result, err = pl.pemAutoEncoding(result)
	return result, nil
}

func (pl *payload[T]) pemAutoEncoding(result []byte) ([]byte, error) {
	if pl.pemType == "" {
		return result, nil
	}
	block := &pem.Block{Type: pl.pemType, Headers: make(map[string]string), Bytes: result}
	for _, customHeader := range pl.customHeaders {
		customHeader(block.Headers)
	}
	output := bytes.NewBuffer(nil)
	err := pem.Encode(output, block)
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func WithData(data any) PayloadOption {
	return func(p *payload[any]) {
		// 处理 data 为 nil 的情况
		if data == nil {
			p.data = nil
			return
		}

		val := reflect.ValueOf(data)
		// 仅当类型为「结构体且非指针」时，转换为指针
		if val.Kind() == reflect.Struct {
			// 创建结构体指针并赋值原始数据
			newPtr := reflect.New(val.Type()) // 生成 *T 类型指针
			newPtr.Elem().Set(val)            // *newPtr = data（将原始结构体值赋给指针指向的对象）
			p.data = newPtr.Interface()       // 转为 any 类型赋值
		} else {
			// 其他类型（指针、基本类型、切片、map 等）直接赋值
			p.data = data
		}
	}
}

func WithLicense(license License) PayloadOption {
	return func(p *payload[any]) {
		WithData(license)(p)
		WithPEM("LICENCE DATA")(p)
	}
}

func WithPEM(pemType string, headers ...CustomPEMHeader) PayloadOption {
	return func(pl *payload[any]) {
		if pemType != "" {
			pl.pemType = pemType
			pl.customHeaders = append(pl.customHeaders, innerHeader(pl))
			for _, pemHeader := range headers {
				pl.customHeaders = append(pl.customHeaders, pemHeader)
			}
		}
	}
}

func WithDeadline(deadline string) PayloadOption {
	return func(p *payload[any]) {
		if deadline != "" {
			p.deadline = deadline
		}
	}
}

func WithVersion(version string) PayloadOption {
	return func(p *payload[any]) {
		if version != "" {
			p.version = version
		}
	}
}
