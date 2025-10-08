package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type OAuth2AccessToken struct {
	ent.Schema
}

func (OAuth2AccessToken) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Comment("令牌ID").Unique(),
		field.String("client_id").Comment("客户端ID").MaxLen(64),
		field.String("subject").Comment("令牌主体（用户）").MaxLen(64).Optional().Nillable(),
		field.Strings("scopes").Comment("授权范围"),
		field.Strings("audience").Comment("受众"),
		field.Time("issued_at").Comment("签发时间").Default(time.Now),
		field.Time("expires_at").Comment("过期时间"),
		field.JSON("metadata", &map[string]interface{}{}).Comment("其他参数").Optional(),
		FieldCreatedAt,
	}
}

func (OAuth2AccessToken) Indexes() []ent.Index {
	return []ent.Index{}
}

func (OAuth2AccessToken) Edges() []ent.Edge {
	return []ent.Edge{
		//edge.To("authorization_code", OAuth2AuthorizationCode.Type).Unique().Annotations(
		//	entsql.OnDelete(entsql.Cascade),
		//),
	}
}

func (OAuth2AccessToken) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "oauth2_access_token"},
		schema.Comment("OAuth2访问令牌"),
	}
}
