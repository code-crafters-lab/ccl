package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().NonNegative(),
		field.String("username").Comment("用户名").MaxLen(32).Unique().NotEmpty(),
		field.String("password").Comment("密码").MaxLen(128).Optional().Nillable(),
		field.String("email").Comment("邮箱").MaxLen(64).Unique().Optional().Nillable(),
		field.Bool("email_verified").Comment("邮箱是否已验证").Default(false),
		// 手机号（可选，用于 MFA 或登录）
		field.String("phone").Comment("手机号").Unique().Optional().Nillable(),
		// 手机号是否验证
		field.Bool("phone_verified").Comment("手机是否已验证").Default(false),
		// 用户状态（枚举：启用/禁用/锁定，默认启用）
		field.Enum("status").Comment("用户状态").Values("active", "inactive", "locked").Default("active"),
		// ABAC 属性：用户核心属性（部门ID、用户ID 等）
		field.JSON("attributes", map[string]interface{}{}).Optional().
			Comment("用户属性，如 {\\\"dept_id\\\": 1001, \\\"user_id\\\": 123}"),
		// 时间戳（自动维护）
		FieldCreatedAt,
		FieldUpdatedAt,
		// 最后登录时间
		field.Time("last_login_at").Comment("最后登录时间").Optional().Nillable(),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email"),      // 高频登录查询字段
		index.Fields("username"),   // 用户名登录查询
		index.Fields("status"),     // 按状态筛选用户（如禁用用户过滤）
		index.Fields("created_at"), // 按注册时间统计/查询
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		//// 1:N 关联：用户 → 凭证（一个用户可有多类凭证：密码、TOTP 等）
		//edge.To("credentials", Credential.Type),
		// M:N 关联：用户 → 角色（通过中间表 UserRole 关联，支持角色过期）
		edge.To("roles", Role.Type).Through("user_roles", UserRole.Type),
		//// 1:N 关联：用户 → 授权码（一个用户可生成多个授权码）
		//edge.To("auth_codes", AuthorizationCode.Type),
		//// 1:N 关联：用户 → 访问令牌
		//edge.To("access_tokens", AccessToken.Type),
		//// 1:N 关联：用户 → 刷新令牌
		//edge.To("refresh_tokens", RefreshToken.Type),
		//// M:N 关联：用户 → 客户端应用（记录用户授权过的应用）
		//edge.To("authorized_clients", ClientApp.Type).Through("user_client_authorizations", UserClientAuthorization.Type),
	}
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		//entsql.Annotation{Table: "Users"},
		schema.Comment("用户信息"),
	}
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		//BaseMixin{},
		//DeleteMixin{},
	}
}
