package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type UserRole struct {
	ent.Schema
}

func (UserRole) Fields() []ent.Field {
	return []ent.Field{
		// 1. 外键 user_id：类型必须与 User 表的 user_id 一致（string）
		field.Int64("user_id").Comment("用户 ID"),
		// 2. 外键 role_id：类型必须与 Role 表的 role_id 一致（string）
		field.Int64("role_id").Comment("角色 ID"),
		// 角色过期时间（可选，支持临时角色）
		field.Time("expired_at").Comment("过期时间").Optional(),
		// 分配角色的操作者（如管理员 ID）
		field.String("assigned_by").Comment("分配角色的操作者").Optional(),
		FieldCreatedAt,
	}
}

func (UserRole) Edges() []ent.Edge {
	return []ent.Edge{
		// 3. 关联 User 表：关键配置
		edge.To("user", User.Type).
			Field("user_id"). // 绑定中间表的 "user_id" 字段
			Unique().         // 确保一个 UserRole 只关联一个 User
			Required(),       // 必选（外键非空）

		// 4. 关联 Role 表：关键配置
		edge.To("role", Role.Type).
			Field("role_id"). // 绑定中间表的 "role_id" 字段
			Unique().         // 确保一个 UserRole 只关联一个 Role
			Required(),       // 必选（外键非空）
	}
}

func (UserRole) Indexes() []ent.Index {
	// 联合唯一：用户+角色唯一（避免重复分配）
	return []ent.Index{
		index.Fields("user_id", "role_id").Unique(),
		// 过期时间索引（用于清理过期角色关联）
		//index.Fields("expired_at").Where(field.NotNull()),
	}
}

func (UserRole) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("用户-角色关联表"),
	}
}
