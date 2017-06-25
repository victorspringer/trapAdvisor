package authenticating

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (s *service) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") == "production" {
			cookieTravellerID, err := r.Cookie("travellerID")
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(401), 401)
				return
			}

			cookieSessionToken, err := r.Cookie("sessionToken")
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(401), 401)
				return
			}

			id, err := strconv.Atoi(cookieTravellerID.Value)
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(401), 401)
				return
			}

			if err = s.ValidateSession(id, cookieSessionToken.Value); err != nil {
				cookieTravellerID.Path = "/"
				cookieTravellerID.MaxAge = -1
				cookieSessionToken.Path = "/"
				cookieSessionToken.MaxAge = -1

				http.SetCookie(w, cookieTravellerID)
				http.SetCookie(w, cookieSessionToken)

				log.Println(err)
				http.Error(w, http.StatusText(401), 401)
				return
			}

			expires := time.Now().Add(30 * 24 * time.Hour)

			cookieTravellerID.Path = "/"
			cookieTravellerID.Expires = expires
			cookieSessionToken.Path = "/"
			cookieSessionToken.Expires = expires

			http.SetCookie(w, cookieTravellerID)
			http.SetCookie(w, cookieSessionToken)
		}

		next(w, r)
	})
}
