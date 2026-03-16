package supervisor

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var supportedSignals = []os.Signal{
	syscall.SIGUSR1,
	syscall.SIGUSR2,
	syscall.SIGHUP,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func Run() error {
	s := &supervisor{}

	if err := s.start(); err != nil {
		return fmt.Errorf("initial start: %w", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, supportedSignals...)

	for sig := range sigCh {
		switch sig {
		case syscall.SIGUSR2:
			log.Println("SIGUSR2: (re)starting ...")
			if err := s.start(); err != nil {
				log.Printf("(re)start failed: %v", err)
			}

		case syscall.SIGUSR1:
			log.Println("SIGUSR1: stopping ...")
			if err := s.stop(); err != nil {
				log.Printf("stop failed: %v", err)
			}

		case syscall.SIGHUP:
			log.Println("SIGHUP: refreshing RU CIDRs and geodata ...")
			if err := s.refresh(); err != nil {
				log.Printf("refresh failed: %v", err)
			}

		case syscall.SIGTERM, syscall.SIGINT:
			log.Println("shutting down ...")
			return s.stop()
		}
	}
	return nil
}
