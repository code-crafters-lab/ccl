package test

import (
	"ccl/db/ent/oauth2client"
	"ccl/db/oauth2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OAuth2ClientAdd(t *testing.T) {
	settings := &oauth2.ClientSettings{
		ResponseTypes: []string{"code", "id_token", "id_token token"},
	}
	c, e := client.OAuth2Client.Create().
		SetID("web").SetAppType("web").SetClientSettings(settings).
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
