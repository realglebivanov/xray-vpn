package main

import (
	"fmt"
	"log"
	"syscall"

	"github.com/realglebivanov/hstd/hstdlib"
	"github.com/realglebivanov/hstd/hstdlib/secret"
	"github.com/realglebivanov/hstd/xrayvpnd/internal/config/store"
	"github.com/spf13/cobra"
)

func newRotateCmd() *cobra.Command {
	return &cobra.Command{
		Use:                "rotate <secret>",
		Short:              "Rotate client UUID based on today's date and a shared secret",
		DisableFlagParsing: true,
		Args:               cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rootSecret, err := hstdlib.ParseHexSecret(args[0])
			if err != nil {
				return fmt.Errorf("secret must be hex: %w", err)
			}
			uuid := secret.GenerateClientUUID(0, rootSecret)
			log.Printf("rotating client_id to %s", uuid)

			if err := store.RotateUUID(uuid); err != nil {
				return err
			}

			fmt.Println("client_id rotated")
			return send(xrayvpndProcess, syscall.SIGUSR2)
		},
	}
}
