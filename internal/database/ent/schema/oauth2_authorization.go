package schema

import (
	"ccl/db/oauth2"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type OAuth2Authorization struct {
	ent.Schema
}

func (OAuth2Authorization) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Comment("授权ID").Unique(),
		field.String("subject").Comment("授权主体（用户）").MaxLen(64).Optional().Nillable(),
		field.String("client_id").Comment("客户端ID").MaxLen(64),
		field.String("response_type").Comment("响应类型").MaxLen(16),
		field.Text("redirect_uri").Comment("重定向地址"),
		field.Strings("scopes").Comment("授权范围"),
		field.String("state").Comment("防CSRF令牌").MaxLen(128).Optional().Nillable(),
		field.String("nonce").Comment("防重放攻击令牌").MaxLen(128).Optional().Nillable(),
		field.String("code_challenge_method").Comment("PKCE").MaxLen(16).Optional().Nillable(),
		field.String("code_challenge").Comment("PKCE").MaxLen(128).Optional().Nillable(),

		field.String("response_mode").Comment("").Optional().Nillable(),
		field.Bool("finished").Default(false).Comment("授权是否已完成"),
		field.Time("auth_time").Comment("授权时间").Optional().Nillable(),
		field.JSON("attributes", &oauth2.AuthorizationAttributes{}).Comment("其他参数"),
		FieldCreatedAt,
	}
}

func (OAuth2Authorization) Indexes() []ent.Index {
	return []ent.Index{}
}

func (OAuth2Authorization) Edges() []ent.Edge {
	return []ent.Edge{
		//edge.To("code", OAuth2AuthorizationCode.Type).Unique().Field("authorization_id"),
	}
}

func (OAuth2Authorization) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "oauth2_authorization"},
		schema.Comment("OAuth2授权信息"),
	}
}
