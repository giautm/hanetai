package hanetai

import (
	"golang.org/x/oauth2"
)

func NewOAuth2Config(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"full"},
		Endpoint: oauth2.Endpoint{
			AuthStyle: oauth2.AuthStyleInParams,
			AuthURL:   "https://oauth.hanet.com/oauth2/authorize",
			TokenURL:  "https://oauth.hanet.com/token",
		},
	}
}
