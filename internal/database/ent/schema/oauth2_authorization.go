package schema

import (
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
		field.String("subject").Comment("授权主体（用户）"),
		field.String("client_id").Comment("客户端ID").MaxLen(64),
		field.String("redirect_uri").Comment("重定向地址"),
		field.String("response_type").Comment("响应类型"),
		field.Strings("scopes").Comment("授权范围"),
		field.String("state").Comment("").MaxLen(128).Optional().Nillable(),
		field.String("nonce").Comment("").MaxLen(128).Optional().Nillable(),
		field.String("code_challenge_method").Comment("").MaxLen(16).Optional().Nillable(),
		field.String("code_challenge").Comment("").MaxLen(128).Optional().Nillable(),
		field.String("response_mode").Comment("").Optional().Nillable(),

		field.Bool("done").Default(false).Comment("授权是否已完成"),
		field.Time("auth_time").Comment("授权时间"),
		field.JSON("attributes", map[string]interface{}{}).Comment("其他参数"),
		FieldCreatedAt,
	}
}

func (OAuth2Authorization) Indexes() []ent.Index {
	return []ent.Index{}
}

func (OAuth2Authorization) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (OAuth2Authorization) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "oauth2_authorization"},
		schema.Comment("OAuth2授权信息"),
	}
}
