package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Tenant holds the schema definition for the Tenant entity.
type Tenant struct {
	ent.Schema
}

// Fields of the Tenant.
func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("名称").NotEmpty().MaxRuneLen(16),
		field.String("description").Comment("描述信息").Optional().MaxLen(64),
	}
}

// Annotations of the Tenant.
func (Tenant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("租户信息"),
	}
}

// Mixin of the Tenant schema.
func (Tenant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		//BaseMixin{},
	}
}

// Policy defines the privacy policy of the User.
//func (Tenant) Policy() ent.Policy {
//	return privacy.Policy{
//		Mutation: privacy.MutationPolicy{
//			privacy.AlwaysDenyRule(),
//		},
//	}
//}
