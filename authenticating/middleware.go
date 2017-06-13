package authenticating

import (
	"net/http"
	"time"
)

func (s *service) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("travellerID")
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		cookie.Path = "/"
		cookie.Expires = time.Now().Add(30 * 24 * time.Hour)
		http.SetCookie(w, cookie)
		next(w, r)
	})
}
