package test

import (
	"ccl/db/ent"
	"ccl/db/ent/oauth2authorization"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OAuth2AuthorizationAdd(t *testing.T) {
	c, e := client.OAuth2Authorization.Create().
		SetClientID("web").SetResponseType("code").SetRedirectURI("http://localhost:3000").
		SetScopes([]string{"read", "write"}).SetState("state").SetNonce("nonce").
		SetCodeChallengeMethod("S256").SetCodeChallenge("code_challenge").
		Save(ctx)

	assert.Nil(t, e)
	assert.NotNil(t, c)
}

func Test_OAuth2AuthorizationQuery(t *testing.T) {
	WithTx(func(tx *ent.Tx) error {

		auth, e1 := client.OAuth2Authorization.Query().Where(oauth2authorization.IDEQ(4)).WithAuthorizationCode().Only(ctx)
		assert.Nil(t, e1)
		assert.NotNil(t, auth)

		code, err := client.OAuth2AuthorizationCode.Create().SetCode("88888").
			SetAuthorizationID(3).Save(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, code)

		return nil
	})

}
