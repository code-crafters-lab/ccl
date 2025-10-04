package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type RoleAuthority struct {
	ent.Schema
}

func (RoleAuthority) Fields() []ent.Field {
	return []ent.Field{
		// 备注（如权限分配原因）
		field.String("remark").Optional(),
		// 时间戳
		field.Time("created_at").Default(time.Now()).Immutable(),
	}
}

func (RoleAuthority) Edges() []ent.Edge {
	return []ent.Edge{
		// 关联角色（必选）
		//edge.From("role", Role.Type).Ref("permissions").Unique().Required(),
		// 关联权限（必选）
		//edge.From("permission", Authority.Type).Ref("roles").Unique().Required(),
	}
}

func (RoleAuthority) Indexes() []ent.Index {
	// 联合唯一：角色+权限唯一（避免重复分配）
	return []ent.Index{
		//index.Fields("role_id", "permission_id").Unique(),
	}
}
