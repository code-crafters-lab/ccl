package test

import (
	"ccl/db/ent/oauth2client"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OAuth2ClientAdd(t *testing.T) {
	c, e := client.OAuth2Client.Create().
		SetID("web").SetAppType("web").
		SetResponseTypes([]string{"code", "id_token", "id_token token"}).
		SetAuthorizationGrantTypes([]string{"authorization_code", "refresh_token"}).
		SetRedirectUris([]string{"http://localhost:5000/auth/callback"}).
		SetAuthenticationMethod("none").
		Save(ctx)

	assert.Nil(t, e)
	assert.NotNil(t, c)
}

func Test_OAuth2ClientQuery(t *testing.T) {
	c, e := client.OAuth2Client.Query().
		Where(oauth2client.IDEQ("web")).
		Only(ctx)

	assert.Nil(t, e)
	assert.NotNil(t, c)
}
