package authenticating

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	fb "github.com/huandu/facebook"
	uuid "github.com/satori/go.uuid"
	"github.com/victorspringer/trapAdvisor/friendship"
	"github.com/victorspringer/trapAdvisor/persistence"
	"github.com/victorspringer/trapAdvisor/traveller"
	"golang.org/x/oauth2"
)

func (s *service) HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(s.config.Endpoint.AuthURL)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parameters := url.Values{}
	parameters.Add("client_id", s.config.ClientID)
	parameters.Add("scope", strings.Join(s.config.Scopes, " "))
	parameters.Add("redirect_uri", s.config.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", s.state)

	u.RawQuery = parameters.Encode()
	url := u.String()

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *service) HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	state := r.FormValue("state")
	if state != s.state {
		err := fmt.Errorf("invalid oauth state, expected '%v', got '%v'", s.state, state)
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	code := r.FormValue("code")
	token, err := s.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := s.config.Client(oauth2.NoContext, token)

	s.session = &fb.Session{
		Version:    "v2.9",
		HttpClient: client,
	}

	param := fb.Params{"access_token": url.QueryEscape(token.AccessToken)}

	trav, err := s.session.Get("/me", param)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(trav)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var t traveller.Traveller
	if err = json.Unmarshal(body, &t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	travRepo := persistence.NewTravellerRepository()

	firstLogin := false
	_, err = travRepo.Find(t.ID)
	if err != nil {
		firstLogin = true
	}

	t.SessionToken = fmt.Sprintf("%v", uuid.NewV4())

	if err = travRepo.Store(&t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if firstLogin {
		friends, err := s.session.Get("/me/friends", fb.Params{"access_token": url.QueryEscape(token.AccessToken), "fields": "id"})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		f := friendship.Friendship{}
		f.TravellerID = t.ID
		fRepo := persistence.NewFriendshipRepository()
		idx := 0
		for friends.Get(fmt.Sprintf("data.%v.id", idx)) != nil {
			id, ok := friends.Get(fmt.Sprintf("data.%v.id", idx)).(string)
			if !ok {
				err = errors.New("invalid user id")
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fID, err := strconv.Atoi(id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			f.FriendID = fID

			if err = fRepo.Store(&f); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			idx++
		}
	}

	expiration := time.Now().Add(30 * 24 * time.Hour)

	cookieTravellerID := http.Cookie{Name: "travellerID", Value: strconv.Itoa(t.ID), Path: "/", Expires: expiration}
	http.SetCookie(w, &cookieTravellerID)

	cookieSessionToken := http.Cookie{Name: "sessionToken", Value: t.SessionToken, Path: "/", Expires: expiration}
	http.SetCookie(w, &cookieSessionToken)

	w.WriteHeader(http.StatusOK)
}

func (s *service) HandleFacebookLogout(w http.ResponseWriter, r *http.Request) {
	cookies := []string{"travellerID", "sessionToken"}

	for _, c := range cookies {
		cookie, err := r.Cookie(c)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(401), 401)
			return
		}

		cookie.Path = "/"
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *service) ValidateSession(id int, sessionToken string) error {
	repo := persistence.NewTravellerRepository()
	if err := repo.FindBySessionToken(id, sessionToken); err != nil {
		return err
	}

	return nil
}
