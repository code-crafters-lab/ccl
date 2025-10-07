package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type OAuth2AuthorizationCode struct {
	ent.Schema
}

func (OAuth2AuthorizationCode) Fields() []ent.Field {
	return []ent.Field{
		field.Int("authorization_id").Comment("授权 ID"),
		field.String("code").Comment("授权码"),
		field.String("state").Comment("防止重放攻击参数").Deprecated("主表已有"),
		field.Bool("is_used").Comment("是否已使用").Default(false),
		field.Time("issued_at").Comment("签发时间").Default(time.Now).Immutable(),
		field.Time("expires_at").Comment("过期时间").Default(func() time.Time {
			return time.Now().Add(5 * time.Minute) // 默认5分钟过期
		}).Immutable(),
		field.JSON("metadata", map[string]interface{}{}).Comment("元数据"),
	}
}

func (OAuth2AuthorizationCode) Indexes() []ent.Index {
	return []ent.Index{}
}

func (OAuth2AuthorizationCode) Edges() []ent.Edge {
	return []ent.Edge{
		//edge.From("authorization", OAuth2Authorization.Type).
		//	Ref("authorization_code").
		//	Unique().Required(),
	}
}

func (OAuth2AuthorizationCode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "oauth2_authorization_code"},
		schema.Comment("OAuth2授权码"),
	}
}
