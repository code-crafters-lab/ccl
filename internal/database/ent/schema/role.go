package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return nil
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return nil
}

// Indexes of the Role.
func (Role) Indexes() []ent.Index {
	return []ent.Index{}
}

// Annotations of the Role.
func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		//entsql.Annotation{Table: "Users"},
		schema.Comment("角色信息"),
	}
}
