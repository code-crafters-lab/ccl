package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

var (
	// FieldCreatedAt 创建时间：默认当前时间，仅初始化时设置
	FieldCreatedAt = field.Time("created_at").Comment("创建时间").
			Default(time.Now).Immutable().Optional().Nillable().
			SchemaType(map[string]string{
			dialect.MySQL: "datetime(3)", // 适配MySQL datetime类型
		},
		)
	// FieldUpdatedAt 更新时间：更新时自动刷新
	FieldUpdatedAt = field.Time("updated_at").Comment("更新时间").
			UpdateDefault(time.Now).Optional().Nillable().
			SchemaType(map[string]string{
			dialect.MySQL: "datetime(3)",
		})
	// FieldDeletedAt 逻辑删除
	FieldDeletedAt = field.Time("deleted_at").Comment("逻辑删除").
			Optional().Nillable().
			SchemaType(map[string]string{
			dialect.MySQL: "datetime(3)",
		})
	// FieldCreatedBy 创建人：可选、不可修改
	FieldCreatedBy = field.String("created_by").Comment("创建人").Optional().Nillable().Immutable().MaxLen(32)
	// FieldUpdatedBy 更新人：可选
	FieldUpdatedBy = field.String("updated_by").Comment("更新人").Optional().Nillable().MaxLen(32)
)

type BaseMixin struct {
	// We embed the `mixin.Schema` to avoid
	// implementing the rest of the methods.
	mixin.Schema
}

// TenantMixin for embedding the tenant info in different schemas.
type TenantMixin struct {
	mixin.Schema
}

// Fields for all schemas that embed TenantMixin.
func (TenantMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tenant_id").Comment("租户 ID").Immutable(),
	}
}

// Edges for all schemas that embed TenantMixin.
func (TenantMixin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tenant", Tenant.Type).Field("tenant_id").Unique().Required().Immutable(),
	}
}

// Policy for all schemas that embed TenantMixin.
func (TenantMixin) Policy() ent.Policy {
	//return rule.FilterTenantRule()
	return nil
}
