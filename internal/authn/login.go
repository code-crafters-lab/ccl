package authn

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type Authentication interface {
	LoginByUsernamePassword(ctx context.Context, username, password, authId string) error
}

type login struct {
	authentication Authentication
	router         chi.Router
	callback       func(context.Context, string) string
}

func NewLogin(authentication Authentication, callback func(context.Context, string) string, issuerInterceptor *op.IssuerInterceptor) *login {
	l := &login{
		authentication: authentication,
		callback:       callback,
	}
	l.createRouter(issuerInterceptor)
	return l
}

func (l *login) createRouter(issuerInterceptor *op.IssuerInterceptor) {
	l.router = chi.NewRouter()
	l.router.Get("/", l.loginHandler)
	l.router.Post("/", issuerInterceptor.HandlerFunc(l.checkLoginHandler))
}

func (l *login) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	// the oidc package will pass the id of the auth request as query parameter
	// we will use this id through the login process and therefore pass it to the login page
	renderLogin(w, r.FormValue(queryAuthRequestID), nil)
}

func renderLogin(w http.ResponseWriter, id string, err error) {
	data := &struct {
		ID    string
		Error string
	}{
		ID:    id,
		Error: errMsg(err),
	}
	err = templates.ExecuteTemplate(w, "login", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (l *login) checkLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	id := r.FormValue("id")
	err = l.authentication.LoginByUsernamePassword(r.Context(), username, password, id)
	if err != nil {
		renderLogin(w, id, err)
		return
	}
	http.Redirect(w, r, l.callback(r.Context(), id), http.StatusFound)
}
