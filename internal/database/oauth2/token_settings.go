package oauth2

import "time"

type TokenSettings struct {
	AuthorizationCodeTimeToLive *time.Duration `json:"authorization_code_time_to_live"` // 授权码过期时间
	AccessTokenTimeToLive       *time.Duration `json:"access_token_time_to_live"`       // 访问令牌过期时间
	AccessTokenFormat           *string        `json:"access_token_format"`             // 访问令牌格式
	DeviceCodeTimeToLive        *time.Duration `json:"device_code_time_to_live"`        // 设备码过期时间
	ReuseRefreshTokens          *bool          `json:"reuse_refresh_tokens"`            // 是否重用刷新令牌
	RefreshTokenTimeToLive      *time.Duration `json:"refresh_token_time_to_live"`      // 刷新令牌过期时间
	IdTokenTimeToLive           *time.Duration `json:"id_token_time_to_live"`           // ID令牌过期时间
	IdTokenSignatureAlgorithm   *string        `json:"id_token_signature_algorithm"`    // ID令牌签名算法
}
