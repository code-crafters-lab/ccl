package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	sessionCookieName  = "ccl_session"
	redirectCookieName = "ccl_redirect_to"
)

// AuthMiddleware 检查用户是否登录，并在未登录时保存重定向 URL
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查用户是否已登录 (这里用一个简单的 cookie 检查作为示例)
		// 在实际应用中，你应该验证 session ID 或 JWT token
		_, err := r.Cookie(sessionCookieName)
		if err == nil {
			// 用户已登录，继续处理请求
			next.ServeHTTP(w, r)
			return
		}

		// 用户未登录，保存当前请求的 URL
		// 注意：不要保存登录和注册页面本身，否则会陷入无限循环
		if !strings.HasPrefix(r.URL.Path, "/login") {
			// 创建一个临时 cookie 来存储重定向 URL
			http.SetCookie(w, &http.Cookie{
				Name:     redirectCookieName,
				Value:    r.URL.Path, // 只保存路径，不包含域名,这里需要完整路径包括参数
				Path:     "/",
				MaxAge:   60 * 10, // 10分钟有效期
				HttpOnly: true,
				Secure:   false, // 在生产环境中应设为 true (HTTPS)
				SameSite: http.SameSiteLaxMode,
			})
		}

		// 将用户重定向到登录页面
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

// SecureRedirect 安全地重定向到本地 URL
// 防止开放重定向攻击 (Open Redirection Attack)
func SecureRedirect(w http.ResponseWriter, r *http.Request, defaultPath string) {
	// 1. 从 cookie 中获取重定向 URL
	redirectCookie, err := r.Cookie(redirectCookieName)
	var redirectPath string

	if err == nil {
		// 2. 验证 URL 的安全性 (必须是本站的本地路径)
		// 解析相对路径 (如 /dashboard) 会得到一个没有 Host 的 URL
		parsedURL, err := url.Parse(redirectCookie.Value)
		if err == nil && parsedURL.Host == "" { // Host 为空表示是本地路径
			redirectPath = redirectCookie.Value
		}
	}

	// 3. 如果没有有效的重定向 URL，则使用默认路径
	if redirectPath == "" {
		redirectPath = defaultPath
	}

	// 4. 清除临时的 redirect_to cookie
	http.SetCookie(w, &http.Cookie{
		Name:   redirectCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1, // 立即删除
	})

	// 5. 执行重定向
	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}
