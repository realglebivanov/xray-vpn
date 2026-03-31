package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/realglebivanov/hstd/hstdlib"
	"github.com/realglebivanov/hstd/xrayconnectord/internal/server"
)

func main() {
	rootSecret := hstdlib.MustEnvHex("SECRET")

	s, err := server.New(rootSecret)
	if err != nil {
		slog.Error("init server", "err", err)
		os.Exit(1)
	}

	go handleSignals(s)

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		slog.Error("listen", "err", err)
		os.Exit(1)
	}
}

func handleSignals(s *server.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	s.Stop()
}
