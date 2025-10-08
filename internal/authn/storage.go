package authn

import (
	"ccl/db/ent"
	"ccl/db/ent/oauth2authorization"
	"ccl/db/ent/oauth2authorizationcode"
	"ccl/db/ent/oauth2client"
	"ccl/db/ent/user"
	"ccl/db/oauth2"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	pwd "github.com/coffee377/autoctl/pkg/security/password"
	"github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type signingKey struct {
	id        string
	algorithm jose.SignatureAlgorithm
	key       *rsa.PrivateKey
}

func (s *signingKey) SignatureAlgorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *signingKey) Key() any {
	return s.key
}

func (s *signingKey) ID() string {
	return s.id
}

type publicKey struct {
	signingKey
}

func (s *publicKey) ID() string {
	return s.id
}

func (s *publicKey) Algorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *publicKey) Use() string {
	return "sig"
}

func (s *publicKey) Key() any {
	return &s.key.PublicKey
}

type storage interface {
	Authentication
	op.Storage
}

type Storage struct {
	storage
	lock       sync.Mutex
	db         *ent.Client
	signingKey signingKey
}

func NewStorage(client *ent.Client) *Storage {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	return &Storage{
		lock: sync.Mutex{},
		db:   client,
		signingKey: signingKey{
			id:        uuid.NewString(),
			algorithm: jose.RS256,
			key:       key,
		},
	}
}

func (s *Storage) Health(ctx context.Context) error {
	//s.db.Ping(ctx)
	//TODO implement me
	return nil
}

func (s *Storage) LoginByUsernamePassword(ctx context.Context, username, password, authId string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	user, err := s.db.User.Query().Select(user.FieldUsername, user.FieldPassword).
		Where(user.UsernameEQ(username), user.StatusEQ("active")).Only(ctx)
	err = errors.Join(err, user.VerifyPassword(pwd.CreateDelegatingPasswordEncoder(), password))
	if err != nil {
		return fmt.Errorf("账号或密码错误")
	}
	if authId != "" {
		aid, e1 := s.getAuthID(authId)
		_, e2 := s.db.OAuth2Authorization.Update().Where(oauth2authorization.IDEQ(aid)).
			SetSubject(user.Username).SetFinished(true).SetAuthTime(time.Now()).Save(ctx)
		return errors.Join(e1, e2)
	}
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
	create := s.db.OAuth2Authorization.Create()
	create.SetScopes(authReq.Scopes)
	create.SetResponseType(string(authReq.ResponseType))
	create.SetClientID(authReq.ClientID)
	create.SetRedirectURI(authReq.RedirectURI)
	create.SetState(authReq.State)
	create.SetNonce(authReq.Nonce)
	create.SetResponseMode(string(authReq.ResponseMode))
	ccm := string(authReq.CodeChallengeMethod)
	create.SetCodeChallengeMethod(ccm)
	create.SetCodeChallenge(authReq.CodeChallenge)
	attributes := oauth2.ConvertAuthRequest2Attributes(authReq)
	create.SetAttributes(attributes)
	authorization, err := create.Save(ctx)
	return authorization, err
}

func (s *Storage) AuthRequestByID(ctx context.Context, authId string) (op.AuthRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	aid, err := s.getAuthID(authId)
	if err != nil {
		return nil, err
	}
	authorization, err := s.db.OAuth2Authorization.Query().Where(oauth2authorization.ID(aid)).Only(ctx)
	return authorization, err
}

func (s *Storage) getAuthID(authId string) (int, error) {
	aid, err := strconv.Atoi(authId)
	if err != nil {
		return 0, fmt.Errorf("aid 解析错误")
	}
	return aid, err
}

func (s *Storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	ac, err := s.db.OAuth2AuthorizationCode.Query().Where(oauth2authorizationcode.CodeEQ(code)).WithAuthorization().Only(ctx)
	if err != nil {
		return nil, err
	}
	return ac.Edges.Authorization, nil
}

func (s *Storage) SaveAuthCode(ctx context.Context, authId string, code string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	aid, err := s.getAuthID(authId)
	if err != nil {
		return err
	}
	codeCreate := s.db.OAuth2AuthorizationCode.Create()
	codeCreate.SetAuthorizationID(aid).SetCode(code)
	// todo 获取 code 过期配置
	//codeCreate.SetExpiresAt(time.Now().Add(time.Hour * 24))
	err = codeCreate.Exec(ctx)
	return err
}

func (s *Storage) DeleteAuthRequest(ctx context.Context, authId string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	aid, err := s.getAuthID(authId)
	if err != nil {
		return err
	}
	err = s.db.OAuth2Authorization.DeleteOneID(aid).Exec(ctx)
	return err
}

func (s *Storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	//TODO implement me
	return "1", time.Now().Add(time.Minute * 5), err
}

func (s *Storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
	//TODO implement me
	return "1", "refresh_token", time.Now().Add(time.Minute * 5), err
}

func (s *Storage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SigningKey(ctx context.Context) (op.SigningKey, error) {
	return &s.signingKey, nil
}

func (s *Storage) SignatureAlgorithms(ctx context.Context) ([]jose.SignatureAlgorithm, error) {
	return []jose.SignatureAlgorithm{s.signingKey.algorithm}, nil
}

func (s *Storage) KeySet(ctx context.Context) ([]op.Key, error) {
	return []op.Key{&publicKey{s.signingKey}}, nil
}

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

func (s *Storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID, clientID string, scopes []string) error {
	return nil
}

func (s *Storage) SetUserinfoFromRequest(ctx context.Context, userinfo *oidc.UserInfo, request op.IDTokenRequest, scopes []string) error {
	return s.setUserinfo(ctx, userinfo, request.GetSubject(), request.GetClientID(), scopes)
}

// setUserinfo sets the info based on the user, scopes and if necessary the clientID
func (s *Storage) setUserinfo(ctx context.Context, userInfo *oidc.UserInfo, userID, clientID string, scopes []string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	user, err := s.db.User.Query().Where(user.UsernameEQ(userID)).Only(ctx)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			userInfo.Subject = strconv.FormatInt(user.ID, 10)
		case oidc.ScopeEmail:
			userInfo.Email = *user.Email
			userInfo.EmailVerified = oidc.Bool(user.EmailVerified)
		case oidc.ScopeProfile:
			userInfo.PreferredUsername = user.Username
			//userInfo.Name = user.FirstName + " " + user.LastName
			//userInfo.FamilyName = user.LastName
			//userInfo.GivenName = user.FirstName
			//userInfo.Locale = oidc.NewLocale(user.PreferredLanguage)
		case oidc.ScopePhone:
			userInfo.PhoneNumber = *user.Phone
			userInfo.PhoneNumberVerified = user.PhoneVerified
			//case CustomScope:
			//	// you can also have a custom scope and assert public or custom claims based on that
			//	userInfo.AppendClaims(CustomClaim, customClaim(clientID))
		}
	}
	return nil
}

func (s *Storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID, subject, origin string) error {
	//TODO implement me
	return fmt.Errorf("not implemented")
}

func (s *Storage) SetIntrospectionFromToken(ctx context.Context, userinfo *oidc.IntrospectionResponse, tokenID, subject, clientID string) error {
	//TODO implement me
	return fmt.Errorf("not implemented")
}

func (s *Storage) GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (map[string]any, error) {
	//TODO implement me
	return nil, fmt.Errorf("not implemented")
}

func (s *Storage) GetKeyByIDAndClientID(ctx context.Context, keyID, clientID string) (*jose.JSONWebKey, error) {
	//TODO implement me
	return nil, fmt.Errorf("not implemented")
}

func (s *Storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	//TODO implement me
	return nil, fmt.Errorf("not implemented")
}
