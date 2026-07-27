# 许可证通用库设计规范（v1.0）

## 1. 概述

- **项目名称**：CCL License Library
- **目标**：跨平台许可证文件格式通用库
- **平台**：Go 核心库
- **位置**：`internal/license`

---

## 2. 文件格式

### 2.1 整体结构

```
┌──────────────────────────────────────────────────────────────┐
│                    License File                               │
├──────────────────────────────────────────────────────────────┤
│  Header (12 bytes)                                          │
│  ┌────────┬────────┬─────────┬─────────┐                   │
│  │ Magic  │ Vendor │ Version │Reserved │                   │
│  │  4B    │   4B   │   2B    │   2B    │                   │
│  └────────┴────────┴─────────┴─────────┘                   │
├──────────────────────────────────────────────────────────────┤
│  TLV Blocks (变长，可嵌套)                                    │
│  ┌────────┬─────────┬────────────────────────────────┐      │
│  │  Type  │   Len   │           Value                │      │
│  │  1B    │   3B    │          (变长)                 │      │
│  └────────┴─────────┴────────────────────────────────┘      │
└──────────────────────────────────────────────────────────────┘
```

### 2.2 Header 详解

| 字段     | 长度 | 说明                                 |
| -------- | ---- | ------------------------------------ |
| Magic    | 4B   | 文件类型标识：`0xCC4C9A9B`           |
| Vendor   | 4B   | 供应商标识（全 0 = 无厂商）          |
| Version  | 2B   | Major(高1B) + Minor(低1B)，当前 v0.1 |
| Reserved | 2B   | 保留字节                             |

### 2.3 TLV 块类型

| Type      | 名称           | 说明           |
| --------- | -------------- | -------------- |
| 0x01      | Payload        | 许可证核心数据 |
| 0x02      | Security       | 安全信息       |
| 0x03      | Signature      | 签名信息       |
| 0x04      | Device Binding | 设备绑定       |
| 0x05      | Extensions     | 扩展块         |
| 0x80-0xFF | Vendor Custom  | 厂商自定义     |

**TLV 编码规则**：

- Type: 1 字节，0x00 保留
- Len: 3 字节，大端序，最大 16MB
- Value: 可嵌套 TLV，支持树形结构

---

## 3. 核心数据结构

### 3.1 Payload（许可证核心数据）

参考 `license.proto` 定义：

```go
type Payload struct {
    LicenseID    string              // 许可证唯一ID
    ProductID    string              // 产品标识
    Mode         Mode                // 授权模式
    Type         Type                // 许可证类型
    IssuedAt     time.Time           // 颁发时间
    ExpiredAt    time.Time           // 过期时间
    Functions    []FunctionPoint      // 功能点列表
    Extensions   map[string]any      // 扩展字段
    Spec         Spec                // 模式特定数据
}
```

**Mode 枚举**：

- `MODE_UNSPECIFIED` (0): 未指定
- `MODE_CLOUD` (1): 云许可
- `MODE_SOFTWARE` (2): 软许可
- `MODE_HARDWARE` (3): 硬件加密许可
- `MODE_BLOCKCHAIN` (4): 区块链许可

**Type 枚举**：

- `TYPE_UNSPECIFIED` (0): 未指定
- `TYPE_INDIVIDUAL` (1): 个人版
- `TYPE_STANDARD` (2): 标准版
- `TYPE_PROFESSIONAL` (3): 专业版
- `TYPE_ENTERPRISE` (4): 企业版
- `TYPE_CUSTOM` (5): 定制版

**FunctionPoint**：

```go
type FunctionPoint struct {
    BitMask     uint64   // 功能位掩码 (2^n)
    Name        string   // 功能名称
    Description string   // 功能描述
    Sort        int64    // 排序
    Group       []string // 分组
}
```

**Spec（模式特定数据）**：

```go
type Spec interface {
    isSpec()
}

// CloudSpec - 云许可
type CloudSpec struct {
    UserID string
}

// SoftwareSpec - 软许可
type SoftwareSpec struct {
    DeviceFingerprints []string // 设备指纹列表
    MaxDeviceCount     uint32   // 最大绑定设备数
}

// HardwareSpec - 硬件加密许可
type HardwareSpec struct {
    USBDongleSerial   string // USB加密狗序列号
    SM4EncryptedData []byte // SM4加密数据
}

// BlockchainSpec - 区块链许可
type BlockchainSpec struct {
    ContractAddress string // 智能合约地址
    TokenID         string // Token ID
    ChainID         string // 链ID
}
```

### 3.2 Security（安全信息）

参考 `license.proto` 的 `SecurityBlock`：

```go
type Security struct {
    SignatureAlgorithmID string // 签名算法标识
    SignatureLength      uint32 // 签名长度
    SignatureData        []byte // 签名值

    EncryptionAlgorithmID  string // 加密算法标识
    EncryptedDataLength   uint32 // 加密数据长度
    EncryptedData         []byte // 加密数据
}
```

**支持的算法**：

| 类别 | 算法                                         |
| ---- | -------------------------------------------- |
| 签名 | SM2-SM3, HMAC-SHA256, RSA-SHA256, ECDSA-P256 |
| 加密 | AES-256-GCM, SM4, ChaCha20-Poly1305          |

### 3.3 Signature（签名）

```go
type Signature struct {
    Algorithm SignatureAlgorithm // 签名算法
    Data      []byte              // 签名值
    PublicKey []byte              // 公钥（可选）
}

enum SignatureAlgorithm {
    SIGNATURE_ALGORITHM_UNSPECIFIED = 0
    SM2_SM3     = 1
    HMAC_SHA256 = 2
    RSA_SHA256  = 3
    ECDSA_P256  = 4
}
```

### 3.4 Device Binding（设备绑定）

```go
type DeviceBinding struct {
    Fingerprint    string   // 硬件指纹（SHA-256）
    Strategy       string   // 绑定策略 (strict/loose)
    BoundDevices   []string // 已绑定设备列表
}
```

**硬件指纹采集**（复用现有）：

- Windows: WMIC
- Linux: /sys/class/dmi/id
- macOS: IOReg

---

## 4. 核心接口定义

```go
type LicenseFile interface {
    Header() Header
    GetBlock(t byte) (TLVBlock, error)
    Blocks() []TLVBlock
    AddBlock(TLVBlock)
    RawBytes() ([]byte, error)
    Verify() error
}

type TLVBlock interface {
    Type() byte
    Value() []byte
    RawBytes() ([]byte, error)
}

type PayloadSerializer interface {
    Serialize(Payload) ([]byte, error)
    Deserialize([]byte) (Payload, error)
}
```

---

## 5. 版本兼容性

- **当前版本**：v0.1 (Major=0, Minor=1)
- **向前兼容**：新版本增加 TLV 块时不影响旧版本解析
- **向后兼容**：解析时忽略未知 Type 块

---

## 6. 实现计划

1. **重构 Header** - 适配新格式（12B + TLV）
2. **新建 TLV 模块** - TLV 块解析/构建/嵌套
3. **重构 Payload** - 整合 Mode/Type/FunctionPoint/Spec
4. **完善 Security** - 参考 proto 的 SecurityBlock
5. **完善 Signature** - 添加算法枚举支持
6. **扩展 Device** - 支持多种绑定策略
7. **编写测试** - 各模块单元测试
8. **性能优化** - 内存分配、并发处理
