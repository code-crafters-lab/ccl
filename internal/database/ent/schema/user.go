package schema

import (
	database "ccl/db"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().NonNegative(),
		field.String("username").Comment("用户名").MaxLen(32).Unique().NotEmpty(),
		field.String("email").Comment("邮箱").MaxLen(64).Unique().Optional(),
		field.Bool("email_verified").Comment("邮箱是否已验证").Nillable().Default(false),
		field.String("password").Comment("密码").MaxLen(128).Optional(),
		// ABAC 属性：用户核心属性（部门ID、用户ID 等）
		field.JSON("attributes", map[string]interface{}{}).
			Optional().
			Comment("用户属性，如 {\\\"dept_id\\\": 1001, \\\"user_id\\\": 123}"),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		// unique index.
		index.Fields("username").Unique(),
		index.Fields("email").Unique(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		//entsql.Annotation{Table: "Users"},
		schema.Comment("用户信息"),
	}
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		database.BaseMixin{},
	}
}
