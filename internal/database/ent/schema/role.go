package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().NonNegative(),
		// 角色唯一标识（如 "admin", "user:read"，用于权限校验）
		field.String("code").Unique().NotEmpty(),
		// 角色名称（展示用，如 "系统管理员"）
		field.String("name").NotEmpty(),
		// 角色描述
		field.String("description").Optional(),
		// 角色状态（启用/禁用）
		field.Bool("enabled").Default(true),
		// 时间戳
		FieldCreatedAt,
		FieldCreatedBy,
		FieldUpdatedAt,
		FieldUpdatedBy,
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		// M:N 关联：角色 → 用户（反向关联 User 的 roles 边）
		edge.From("users", User.Type).Ref("roles").
			Through("user_roles", UserRole.Type),

		// M:N 关联：角色 → 权限（通过中间表 RolePermission 关联）
		//edge.To("permissions", Permission.Type).Through("role_permissions", RolePermission.Type),
	}
}

// Indexes of the Role.
func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code"),    // 权限校验时高频查询角色编码
		index.Fields("enabled"), // 筛选启用的角色
	}
}

// Annotations of the Role.
func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		//entsql.Annotation{Table: "Users"},
		schema.Comment("角色信息"),
	}
}
