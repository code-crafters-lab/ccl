package oauth2

import "github.com/zitadel/oidc/v3/pkg/oidc"

type AuthorizationAttributes struct {
	Audience    []string                 `json:"audience,omitempty"`
	Display     oidc.Display             `json:"display,omitempty"`
	Prompt      oidc.SpaceDelimitedArray `json:"prompt,omitempty"`
	MaxAge      *uint                    `json:"max_age,omitempty"`
	UILocales   oidc.Locales             `json:"ui_locales,omitempty"`
	IDTokenHint string                   `json:"id_token_hint,omitempty"`
	LoginHint   string                   `json:"login_hint,omitempty"`
	ACRValues   oidc.SpaceDelimitedArray `json:"acr_values,omitempty"`
}

func ConvertAuthRequest2Attributes(request *oidc.AuthRequest) *AuthorizationAttributes {
	attributes := &AuthorizationAttributes{
		Audience:    nil,
		Display:     request.Display,
		Prompt:      request.Prompt,
		MaxAge:      request.MaxAge,
		UILocales:   request.UILocales,
		IDTokenHint: request.IDTokenHint,
		LoginHint:   request.LoginHint,
		ACRValues:   request.ACRValues,
	}
	return attributes
}
