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
	"github.com/victorspringer/trapAdvisor/friendship"
	"github.com/victorspringer/trapAdvisor/persistence"
	"github.com/victorspringer/trapAdvisor/traveller"
	"golang.org/x/oauth2"
)

func (s *service) HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(s.config.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
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
		err := fmt.Errorf("invalid oauth state, expected '%s', got '%s'", s.state, state)
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	code := r.FormValue("code")
	token, err := s.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
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
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	body, err := json.Marshal(trav)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	var t traveller.Traveller
	if err = json.Unmarshal(body, &t); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	travRepo := persistence.NewTravellerRepository()

	firstLogin := false
	_, err = travRepo.Find(t.ID)
	if err != nil {
		firstLogin = true
	}

	if err = travRepo.Store(&t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	if firstLogin {
		friends, err := s.session.Get("/me/friends", fb.Params{"access_token": url.QueryEscape(token.AccessToken), "fields": "id"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				log.Fatal(err)
			}
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
				w.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(w).Encode(err); err != nil {
					log.Fatal(err)
				}
				return
			}

			fID, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(w).Encode(err); err != nil {
					log.Fatal(err)
				}
				return
			}

			f.FriendID = fID

			if err = fRepo.Store(&f); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(w).Encode(err); err != nil {
					log.Fatal(err)
				}
				return
			}

			idx++
		}
	}

	expiration := time.Now().Add(30 * 24 * time.Hour)
	cookie := http.Cookie{Name: "travellerID", Value: strconv.Itoa(t.ID), Path: "/", Expires: expiration}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

func (s *service) HandleFacebookLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("travellerID")
	if err != nil {
		http.Error(w, http.StatusText(401), 401)
		return
	}

	cookie.Path = "/"
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}
