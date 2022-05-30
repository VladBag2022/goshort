package server

import (
	"github.com/VladBag2022/goshort/internal/misc"
	"net/http"
)


func authenticationMiddleware(s Server) func(next http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
			cookie, err := r.Cookie(s.config.AuthCookieName)
			setNewCookie := false
			if err == nil {
				validCookie, userID, _ := misc.Verify(s.config.AuthCookieKey, cookie.Value)

				if !validCookie {
					setNewCookie = true
				} else {
					_, err = s.repository.ShortenedList(r.Context(), userID)
					if err != nil {
						setNewCookie = true
					}
				}

			} else if err == http.ErrNoCookie {
				setNewCookie = true

			} else {

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if setNewCookie {
				userID, err := s.repository.Register(r.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				cookie = &http.Cookie{
					Name:  s.config.AuthCookieName,
					Value: misc.Sign(s.config.AuthCookieKey, userID),
				}
				http.SetCookie(w, cookie)
			}

			next.ServeHTTP(w, r)
		})
	}
}