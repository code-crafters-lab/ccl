package schema

import (
	"ccl/db/oauth2"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type OAuth2Client struct {
	ent.Schema
}

func (OAuth2Client) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Comment("客户端 id").MaxLen(64).Unique().Immutable(),
		field.String("secret").Comment("客户端密钥").MaxLen(128).Sensitive().Optional().Nillable(),
		field.Strings("redirect_uris").Comment("重定向地址").Optional(),
		field.Strings("postLogoutRedirectUris").Comment("退出后重定向地址").Optional(),
		//{"value": "1", "code": "NONE", "name": "公开客户端"},
		//{"value": "2", "code": "CLIENT_SECRET_BASIC", "name": "Basic认证"},
		//{"value": "4", "code": "CLIENT_SECRET_POST", "name": "POST认证"},
		//{"value": "8", "code": "CLIENT_SECRET_JWT", "name": "客户端密钥签证JWT"},
		//{"value": "16", "code": "PRIVATE_KEY_JWT", "name": "私钥签证JWT"}
		field.String("authentication_method").MaxLen(24).Comment("客户端身份验证方法"),
		field.Enum("app_type").Comment("应用类型").
			NamedValues("0", "web", "1", "user_agent", "2", "native"),
		field.Strings("authorization_grant_types").Comment("授权类型"),
		field.Bool("is_dev").Comment("开发模式").Default(false),
		field.Bool("claims_assertion").Comment("idTokenUserinfoClaimsAssertion").Default(false),
		field.Int64("time_offset").Comment("时间偏移（s）").Optional().Nillable(),
		field.JSON("client_settings", &oauth2.ClientSettings{}).Comment("客户端配置").Optional(),
		field.JSON("token_settings", &oauth2.TokenSettings{}).Comment("令牌设置").Optional(),
		FieldCreatedAt,
	}
}

func (OAuth2Client) Indexes() []ent.Index {
	return []ent.Index{}
}

func (OAuth2Client) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (OAuth2Client) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "oauth2_client"},
		schema.Comment("OAuth2客户端"),
	}
}

func (OAuth2Client) Hooks() []ent.Hook {
	return []ent.Hook{}
}
