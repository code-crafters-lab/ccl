package authn

import (
	"ccl/db/ent"
	"ccl/db/ent/oauth2client"
	"context"
	"fmt"
	"sync"

	oidcstorage "github.com/zitadel/oidc/v3/example/server/storage"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type Storage struct {
	storage
	lock sync.Mutex
	db   *ent.Client
	S    *oidcstorage.Storage
}

func NewStorage(client *ent.Client) op.Storage {
	store := oidcstorage.NewUserStore("http://localhost")
	return &Storage{
		lock: sync.Mutex{},
		db:   client,
		S:    oidcstorage.NewStorage(store),
	}
}

func (s *Storage) Health(ctx context.Context) error {
	//s.db.Ping(ctx)
	//TODO implement me
	return nil
}

func (s *Storage) CreateAuthRequest(ctx context.Context, authReq *oidc.AuthRequest, uid string) (op.AuthRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(authReq.Prompt) == 1 && authReq.Prompt[0] == "none" {
		// With prompt=none, there is no way for the user to log in
		// so return error right away.
		return nil, oidc.ErrLoginRequired()
	}
	// 数据库存储
	//authorization, err := s.db.OAuth2Authorization.Create().SetScopes(authReq.Scopes).Save(ctx)
	//return authorization, err
	return nil, nil
}

func (s *Storage) AuthRequestByID(ctx context.Context, authId string) (op.AuthRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SaveAuthCode(ctx context.Context, authId string, code string) error {
	//TODO implement me
	panic("implement me")
}

//func (s *Storage) DeleteAuthRequest(ctx context.Context, authId string) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) SigningKey(ctx context.Context) (op.SigningKey, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) SignatureAlgorithms(ctx context.Context) ([]jose.SignatureAlgorithm, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) KeySet(ctx context.Context) ([]op.Key, error) {
//	//TODO implement me
//	panic("implement me")
//}

func (s *Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// TODO 缓存处理
	client, err := s.db.OAuth2Client.Query().Where(oauth2client.IDEQ(clientID)).Only(ctx)
	return client, err
}

func (s *Storage) AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	client, err := s.db.OAuth2Client.Query().Where(oauth2client.IDEQ(clientID)).Only(ctx)
	if err != nil {
		return err
	}
	if *client.Secret != clientSecret {
		return fmt.Errorf("invalid secret")
	}
	return nil
}

//func (s *Storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID, clientID string, scopes []string) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID, subject, origin string) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) SetIntrospectionFromToken(ctx context.Context, userinfo *oidc.IntrospectionResponse, tokenID, subject, clientID string) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (map[string]any, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) GetKeyByIDAndClientID(ctx context.Context, keyID, clientID string) (*jose.JSONWebKey, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Storage) CheckUsernamePassword(username, password, id string) error {
//	//TODO implement me
//	panic("implement me")
//}
