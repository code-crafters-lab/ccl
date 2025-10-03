package database

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// -------------------------------------------------
// Mixin definition

// BaseMixin implements the ent.Mixin for sharing
// time fields with package schemas.
type BaseMixin struct {
	// We embed the `mixin.Schema` to avoid
	// implementing the rest of the methods.
	mixin.Schema
}

func (BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		// 创建时间：默认当前时间，仅初始化时设置
		field.Time("created_at").Comment("创建时间").
			Default(time.Now).
			Immutable().                  // 不可修改
			SchemaType(map[string]string{ // 适配MySQL datetime类型
				dialect.MySQL: "datetime(3)",
			}),
		// 创建人：可选、不可修改
		field.String("created_by").Comment("创建人").MaxLen(32).Optional().Immutable(),
		// 更新时间：更新时自动刷新
		field.Time("updated_at").Comment("更新时间").Optional().
			UpdateDefault(time.Now).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime(3)",
			}),
		// 更新人：可选
		field.String("updated_by").Comment("更新人").MaxLen(32).Optional(),
	}
}
