package authn

import (
	"ccl/db/ent"
	"context"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"syscall"
	"time"

	"ccl/authn/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/run"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"
)

const (
	pathLoggedOut = "/logout"
)

var (
	key = sha256.Sum256([]byte("test"))
)

type AuthorizationServer interface {
	Authentication
	OIDC() (op.OpenIDProvider, error)
	RegisterRouter() (chi.Router, error)
	Run() error
}

type authorizationServer struct {
	logger       *slog.Logger
	storage      storage
	extraOptions []op.Option
}

func (a *authorizationServer) LoginByUsernamePassword(ctx context.Context, username, password, authId string) error {
	return a.storage.LoginByUsernamePassword(ctx, username, password, authId)
}

func (a *authorizationServer) OIDC() (op.OpenIDProvider, error) {
	config := &op.Config{
		CryptoKey: key,
		// will be used if the end_session endpoint is called without a post_logout_redirect_uri
		DefaultLogoutRedirectURI: pathLoggedOut,
		// enables code_challenge_method S256 for PKCE (and therefore PKCE in general)
		CodeMethodS256: true,
		// enables additional client_id/client_secret authentication by form post (not only HTTP Basic Auth)
		AuthMethodPost: true,
		// enables additional authentication by using private_key_jwt
		AuthMethodPrivateKeyJWT: true,
		// enables refresh_token grant use
		GrantTypeRefreshToken: true,
		// enables use of the `request` Object parameter
		RequestObjectSupported: true,
		// this example has only static texts (in English), so we'll set the here accordingly
		SupportedUILocales: []language.Tag{language.English},
		DeviceAuthorization: op.DeviceAuthorizationConfig{
			Lifetime:     5 * time.Minute,
			PollInterval: 5 * time.Second,
			UserFormPath: "/device",
			UserCode:     op.UserCodeBase20,
		},
	}
	handler, err := op.NewProvider(config, a.storage, op.IssuerFromHost(""),
		append([]op.Option{
			//we must explicitly allow the use of the http issuer
			op.WithAllowInsecure(),
			// as an example on how to customize an endpoint this will change the authorization_endpoint from /authorize to /auth
			op.WithCustomIntrospectionEndpoint(op.NewEndpoint("oauth2/introspect")),
			op.WithCustomEndpoints(
				op.NewEndpoint("oauth2/authorize"),
				op.NewEndpoint("oauth2/token"),
				op.NewEndpoint("oauth2/userinfo"),
				op.NewEndpoint("oauth2/revoke"),
				op.NewEndpoint("oauth2/end_session"),
				op.NewEndpoint("oauth2/jwks"),
			),
			// Pass our logger to the OP
			op.WithLogger(a.logger.WithGroup("op")),
		}, a.extraOptions...)...,
	)
	return handler, err
}

func (a *authorizationServer) RegisterRouter() (chi.Router, error) {
	provider, err := a.OIDC()
	if err != nil {
		a.logger.Error("failed to create oidc provider", "error", err)
		return nil, err
	}

	router := chi.NewRouter()

	// --- 公开路由 ---
	router.Group(func(r chi.Router) {
		l := NewLogin(a, op.AuthCallbackURL(provider), op.NewIssuerInterceptor(provider.IssuerFromRequest))
		router.Mount("/login", http.StripPrefix("/login", l.router))
	})

	router.Get("/auth/callback/test", codeHandler)

	// --- 受保护路由 ---
	router.Group(func(r chi.Router) {
		//r.Use(middleware.AuthMiddleware) // 应用认证中间件
		r.Get("/dashboard", dashboardHandler)
		r.Get("/profile", profileHandler)
		r.Get("/logout", logoutHandler)

		r.Route("/device", func(r chi.Router) {
			a.logger.Warn("registering device auth")
			//registerDeviceAuth(storage, r)
		})

		handler := http.Handler(provider)

		// we register the http handler of the OP on the root, so that the discovery endpoint (/.well-known/openid-configuration)
		// is served on the correct path
		//
		// if your issuer ends with a path (e.g. http://localhost:9998/custom/path/),
		// then you would have to set the path prefix (/custom/path/)
		r.Mount("/", handler)
	})

	router.Get("/about", func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte("about"))
	})

	return router, nil
}

func (a *authorizationServer) Run() error {
	g := &run.Group{}
	addHttp(g, a)
	addExtra(g)
	return g.Run()
}

func addHttp(g *run.Group, a *authorizationServer) {
	var httpSrv *http.Server
	g.Add(func() error {
		router, err := a.RegisterRouter()
		if err != nil {
			return err
		}
		httpSrv = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", "", 80),
			Handler: router,
		}
		return httpSrv.ListenAndServe()
	}, func(err error) {
		if err := httpSrv.Close(); err != nil {
			a.logger.Error("failed to stop web server: %v", err)
		}
	})
}

func addExtra(g *run.Group) {
	g.Add(run.SignalHandler(context.TODO(), syscall.SIGINT, syscall.SIGTERM))
}

func NewAuthorizationServer(client *ent.Client) AuthorizationServer {
	server := &authorizationServer{
		logger:  slog.Default(),
		storage: NewStorage(client),
	}
	return server
}

// 简单的模板渲染函数
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, _ := template.New("").Parse(tmpl)
	_ = t.Execute(w, data)
}

// --- 处理器 ---

// 登录页面处理器
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	html := `
		<html>
		<body>
			<h1>Login</h1>
			<form method="post" action="/login">
				Username: <input type="text" name="username"><br>
				Password: <input type="password" name="password"><br>
				<input type="submit" value="Login">
			</form>
		</body>
		</html>
	`
	renderTemplate(w, html, nil)
}

// 登录处理处理器 (处理 POST 请求)
func loginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	// 实际应用中应从表单获取并验证
	username := r.FormValue("username")
	password := r.FormValue("password")

	// 简单的验证逻辑
	if username == "admin" && password == "123456" {
		// 登录成功，创建一个 session cookie
		// 在实际应用中，你应该生成一个安全的、唯一的 session ID
		http.SetCookie(w, &http.Cookie{
			Name:     "user_session",
			Value:    "some-secure-session-id",
			Path:     "/",
			MaxAge:   60 * 60 * 24, // 24小时有效期
			HttpOnly: true,
			Secure:   false, // 在生产环境中应设为 true (HTTPS)
			SameSite: http.SameSiteLaxMode,
		})

		// 使用我们创建的辅助函数进行安全重定向
		// 如果没有 redirect_to cookie，将默认重定向到 /dashboard
		middleware.SecureRedirect(w, r, "/dashboard")
		return
	}

	// 登录失败
	html := `
		<html>
		<body>
			<h1>Login Failed</h1>
			<p>Invalid username or password.</p>
			<a href="/login">Try again</a>
		</body>
		</html>
	`
	renderTemplate(w, html, nil)
}

func codeHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	state := r.FormValue("state")
	log.Printf("code: %s, state: %s", code, state)

	// 1. 准备表单数据
	formData := url.Values{
		"client_id":     {"web"},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"code_verifier": {"BUn8vDyURwERA1n~LtuC7CbGGgsvJp_CSAs-SvRbIZPGXkqFeXJC29Lc.TDZqIrJ"},
		"redirect_uri":  {"http://localhost/auth/callback/test"},
	}
	// 2. 发送 POST 请求
	u := fmt.Sprintf("%s://%s/oauth2/token", "http", r.Host)
	resp, err := http.PostForm(u, formData)
	if err != nil {
		log.Fatal(err)
	}
	// 3. 确保响应体被关闭
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	// 4. 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", resp.StatusCode)
	}

	// 5. 读取并打印响应体
	body, err := io.ReadAll(resp.Body)
	log.Printf("Response body: %s", body)
}

// 登出处理器
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// 清除 session cookie
	//_, _ = w.Write([]byte("signed out successfully"))
	// todo 跳转回统一认证中心登录页
	http.SetCookie(w, &http.Cookie{
		Name:   "user_session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// 受保护的页面
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	html := `
		<html>
		<body>
			<h1>Welcome to Dashboard!</h1>
			<p><a href="/.well-known/openid-configuration">well-known</a></p>
			<p><a href="/oauth2/authorize?client_id=web&response_type=code&scope=openid&redirect_uri=http://localhost/auth/callback/test&state=123456&code_challenge_method=S256&code_challenge=uh5zryRg8UqzGuXk7ao9V0Smo34i7icPBrubVWkDgEw
">SSO</a></p>
			<p><a href="/profile">Go to Profile</a></p>
			<p><a href="/logout">Logout</a></p>
		</body>
		</html>
	`
	renderTemplate(w, html, nil)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	html := `
		<html>
		<body>
			<h1>Your Profile</h1>
			<p><a href="/dashboard">Back to Dashboard</a></p>
			<p><a href="/logout">Logout</a></p>
		</body>
		</html>
	`
	renderTemplate(w, html, nil)
}
