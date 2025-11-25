package license

type License interface {
	//Header() Header
	// ID 1. 基础信息
	ID() string // 唯一标识
	//Type() string         // 许可证类型（如试用、正式、企业版）
	//IssuedAt() time.Time  // 颁发时间
	//ExpiredAt() time.Time // 过期时间
	//
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
