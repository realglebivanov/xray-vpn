package server

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/realglebivanov/hstd/hstdlib/secret"
	"github.com/realglebivanov/hstd/xrayconnectord/internal/client"
	"github.com/realglebivanov/hstd/xrayconnectord/internal/link"
)

func (s *Server) HandleSubReq(w http.ResponseWriter, r *http.Request) {
	l, httpErr := s.validate(r)
	if httpErr != nil {
		http.Error(w, httpErr.reason, httpErr.code)
		return
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	if err := s.db.TrackIP(l, host); err != nil {
		log.Printf("track ip: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	uuid := secret.GenerateClientUUID(l.Index, s.rootSecret)
	configs := client.BuildConfigs(uuid, s.serverConfigs)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("profile-update-interval", "1")
	w.Header().Set("profile-title", "base64:"+base64.StdEncoding.EncodeToString([]byte("hstd")))

	if err := json.NewEncoder(w).Encode(configs); err != nil {
		log.Printf("encode response: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

type httpError struct {
	reason string
	code   int
}

func (s *Server) validate(r *http.Request) (*link.Link, *httpError) {
	l, httpErr := s.buildLink(r)
	if httpErr != nil {
		return nil, httpErr
	}

	isEnabled, err := s.db.IsEnabled(l)
	if err != nil {
		log.Printf("check link enabled: %v", err)
		return nil, &httpError{"internal error", http.StatusInternalServerError}
	}
	if !isEnabled {
		log.Printf("link disabled: %d", l.Index)
		return nil, &httpError{"bad request", http.StatusBadRequest}
	}
	return l, nil
}

func (s *Server) buildLink(r *http.Request) (*link.Link, *httpError) {
	src := r.URL.Path[1:]

	if src == s.legacySubPath {
		return link.New(0, s.rootSecret), nil
	}

	l, err := link.Unmarshal(src)
	if err != nil {
		log.Printf("decode input link: %v", err)
		return nil, &httpError{"bad request", http.StatusBadRequest}
	}

	if !l.IsValid(s.rootSecret) {
		return nil, &httpError{"bad request", http.StatusBadRequest}
	}

	return l, nil
}
