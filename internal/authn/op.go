package authn

import "github.com/zitadel/oidc/v3/pkg/op"

type storage interface {
	Authenticate
	op.Storage
}
