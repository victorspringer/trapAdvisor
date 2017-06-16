package authenticating

import (
	"fmt"
	"net/http"
	"os"

	fb "github.com/huandu/facebook"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

const (
	clientID     = "CLIENT_ID"
	clientSecret = "CLIENT_SECRET"
)

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
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  getDomain() + "/auth_callback",
			Scopes:       []string{"public_profile", "user_friends"},
			Endpoint:     facebook.Endpoint,
		},
		state: fmt.Sprintf("%s", uuid.NewV4()),
	}
}

func getDomain() string {
	if os.Getenv("DOMAIN") != "" {
		return os.Getenv("DOMAIN")
	}
	return "http://localhost:8080"
}
