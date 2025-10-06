package test

import (
	"ccl/db/ent/role"
	"ccl/db/ent/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_User_Add(t *testing.T) {
	err = client.User.Create().
		SetUsername("test").SetPassword("123456").
		SetEmail("test@dingtalk.com").Exec(ctx)
	assert.Nil(t, err)
}

func Test_User_Update(t *testing.T) {
	id, err2 := client.Role.Query().Where(role.CodeEQ("admin")).OnlyID(ctx)
	assert.Nil(t, err2)
	u, err := client.User.Update().Where(user.UsernameEQ("coffee377")).
		SetEmail("coffee377@dingtalk.com").AddRoleIDs(id).
		Save(ctx)
	assert.Equal(t, 1, u)
	assert.Nil(t, err)
}

func Test_User_AddRole(t *testing.T) {
	tx, err = client.Tx(ctx)
	assert.Nil(t, err)
	only, err2 := tx.User.Query().Where(user.UsernameEQ("coffee377")).Only(ctx)
	assert.Nil(t, err2)
	assert.NotNil(t, only)

	err2 = tx.Role.Create().SetCode("admin").SetName("管理员").
		AddUsers(only).Exec(ctx)
	assert.Nil(t, err2)

	err2 = tx.Commit()
	assert.Nil(t, err2)

}

func Test_User_QueryWithRole(t *testing.T) {
	user, err2 := client.User.Query().Where(user.UsernameEQ("coffee377")).
		WithRoles().Only(ctx)
	assert.Nil(t, err2)
	assert.NotNil(t, user)
	assert.NotNil(t, user.Edges.Roles)

}

func Test_User_Delete(t *testing.T) {
	//u, e := client.User.Query().Where(user.UsernameEQ("coffee377")).WithUserRoles().Only(ctx)
	//assert.Nil(t, e)
	//assert.NotNil(t, u)
	//t2, e := client.Tx(ctx)
	//assert.Nil(t, e)
	//exec, e := t2.UserRole.Delete().Where(userrole.HasUserWith(user.UsernameEQ("coffee377"))).Exec(ctx)
	//assert.Nil(t, e)
	//assert.Equal(t, 1, exec)
	//u.Edges.UserRoles
	i, err := client.User.Delete().Where(user.UsernameEQ("coffee377")).Exec(ctx)
	assert.Equal(t, 1, i)
	assert.Nil(t, err)
	//e = t2.Commit()
	//assert.Nil(t, e)

}
