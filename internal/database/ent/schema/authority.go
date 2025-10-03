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
		field.String("resource_id").Comment("资源ID").MaxLen(64).Optional(),
		field.String("parent_resource_id").Comment("父资源ID").MaxLen(64).Optional(),
		field.String("name").Comment("资源名称").MaxLen(32),
		field.Int("sort").Comment("资源排序").Default(1),
		field.String("remark").Comment("资源备注").MaxLen(128).Optional(),

		// todo 补充创建信息

		//DELETE_FLAG   tinyint(1)  default 0  not null comment '数据状态：1_已删除；0_未删除'
	}
}

// Indexes of the Authority.
func (Authority) Indexes() []ent.Index {
	return []ent.Index{
		// unique index.
		index.Fields("resource_type", "resource_id").Unique(),
		index.Fields("parent_resource_id"),
	}
}

// Edges of the Authority.
func (Authority) Edges() []ent.Edge {
	return nil
}

// Annotations of the Authority.
func (Authority) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("资源权限"),
	}
}

// Mixin of the Authority.
func (Authority) Mixin() []ent.Mixin {
	return []ent.Mixin{
		//database.BaseMixin{},
	}
}
