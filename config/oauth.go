package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     "778838656131-jfnap1huoa7igvob44b1159gg0e2q99e.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}
