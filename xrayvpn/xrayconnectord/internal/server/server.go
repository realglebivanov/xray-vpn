package server

import (
	"fmt"
	"log"
	"net"

	"github.com/realglebivanov/hstd/hstdlib"
	"github.com/realglebivanov/hstd/xrayconnectord/internal/client"
	"github.com/realglebivanov/hstd/xrayconnectord/internal/db"
)

type Server struct {
	db            *db.DB
	legacySubPath string
	rootSecret    []byte
	adminCIDR     *net.IPNet
	serverConfigs []*client.ServerConfig
}

func New(rootSecret []byte) (*Server, error) {
	db, err := db.Open()
	if err != nil {
		return nil, fmt.Errorf("open database: %v", err)
	}

	_, adminCIDR, err := net.ParseCIDR(hstdlib.MustEnv("ADMIN_CIDR"))
	if err != nil {
		return nil, fmt.Errorf("parse ADMIN_CIDR: %v", err)
	}

	return &Server{
		db:            db,
		legacySubPath: hstdlib.MustEnv("SUB_PATH"),
		rootSecret:    rootSecret,
		adminCIDR:     adminCIDR,
		serverConfigs: []*client.ServerConfig{{
			Remark:     "Обычный ВПН",
			Host:       hstdlib.MustEnv("SERVER_HOST"),
			RealityPbk: hstdlib.MustEnv("REALITY_PBK"),
			RealitySni: hstdlib.MustEnv("REALITY_SNI"),
			RealitySid: hstdlib.MustEnv("REALITY_SID"),
		}, {
			Remark:     "Обход белых списков",
			Host:       hstdlib.MustEnv("PROXY_HOST"),
			RealityPbk: hstdlib.MustEnv("REALITY_PBK"),
			RealitySni: hstdlib.MustEnv("REALITY_SNI"),
			RealitySid: hstdlib.MustEnv("REALITY_SID"),
		}}}, nil
}

func (s *Server) Close() {
	if err := s.db.Close(); err != nil {
		log.Printf("close db: %v", err)
	}
}
