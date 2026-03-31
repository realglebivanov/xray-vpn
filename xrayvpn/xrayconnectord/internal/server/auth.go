package server

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) basicAuth(w http.ResponseWriter, r *http.Request) bool {
	sesh, _ := s.auth.sessions.Get(r, "hstd#xrayconnectord#subsrv")

	if _, ok := sesh.Values["id"]; ok {
		return true
	}

	user, pass, ok := r.BasicAuth()
	if !ok || user != s.auth.adminUser {
		return false
	}

	if bcrypt.CompareHashAndPassword([]byte(s.auth.adminPasswordHash), []byte(pass)) != nil {
		return false
	}

	sesh.Options = s.auth.sessionOpts
	sesh.Values["id"] = true

	return sesh.Save(r, w) == nil
}
