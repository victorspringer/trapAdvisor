package authenticating

import (
	"net/http"

	uuid "github.com/satori/go.uuid"

	"fmt"

	fb "github.com/huandu/facebook"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

const domain = "http://localhost:8080"

// Service is the interface that provides Facebook authentication methods.
type Service interface {
	HandleFacebookLogin(http.ResponseWriter, *http.Request)
	HandleFacebookCallback(http.ResponseWriter, *http.Request)
	HandleFacebookLogout(http.ResponseWriter, *http.Request)
	AuthMiddleware(http.HandlerFunc) http.HandlerFunc
}

type service struct {
	config  *oauth2.Config
	state   string
	session *fb.Session
}

// NewService creates an authentication service.
func NewService() Service {
	return &service{
		config: &oauth2.Config{
			ClientID:     "CLIENT_ID",
			ClientSecret: "CLIENT_SECRET",
			RedirectURL:  domain + "/auth_callback",
			Scopes:       []string{"public_profile", "user_friends"},
			Endpoint:     facebook.Endpoint,
		},
		state: fmt.Sprintf("%s", uuid.NewV4()),
	}
}
