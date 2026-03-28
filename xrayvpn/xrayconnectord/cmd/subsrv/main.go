package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/realglebivanov/hstd/hstdlib"
	"github.com/realglebivanov/hstd/xrayconnectord/internal/server"
)

func main() {
	rootSecret := hstdlib.MustEnvHex("SECRET")

	s, err := server.New(rootSecret)
	if err != nil {
		log.Fatalf("init server: %v", err)
	}
	defer s.Close()

	http.HandleFunc("/admin/", s.HandleAdminReq)
	http.HandleFunc("/", s.HandleSubReq)

	credsDir := hstdlib.MustEnv("CREDENTIALS_DIRECTORY")
	certFile := filepath.Join(credsDir, "tls_cert")
	keyFile := filepath.Join(credsDir, "tls_key")

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServeTLS(":8080", certFile, keyFile, nil))
}
