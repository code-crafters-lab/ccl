package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Authority holds the schema definition for the Authority entity.
type Authority struct {
	ent.Schema
}

// Fields of the Authority.
func (Authority) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("resource_type").Comment("资源类型").Values("1", "2"),
		field.String("resource_id").Comment("资源ID").MaxLen(64).NotEmpty(),
		field.String("parent_resource_id").Comment("父资源ID").MaxLen(64).Optional(),
		field.String("name").Comment("资源名称").NotEmpty().MaxLen(32),
		field.Int("sort").Comment("资源排序").Default(1),
		field.String("description").Comment("权限描述").MaxLen(128).Optional(),
		FieldCreatedAt,
		FieldUpdatedAt,
	}
}

// Edges of the Authority.
func (Authority) Edges() []ent.Edge {
	return []ent.Edge{
		// M:N 关联：权限 → 角色（反向关联 Role 的 permissions 边）
		//edge.From("roles", Role.Type).Ref("permissions").Through("role_permissions", RolePermission.Type),
	}
}

// Indexes of the Authority.
func (Authority) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("resource_type", "resource_id").Unique(),
		index.Fields("parent_resource_id"),
	}
}

// Annotations of the Authority.
func (Authority) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("资源权限"),
	}
}
