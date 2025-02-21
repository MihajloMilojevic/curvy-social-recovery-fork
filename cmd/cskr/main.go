package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/0x3327/curvy-social-recovery/internal/cskr/commands/recoverCmd"
	"github.com/0x3327/curvy-social-recovery/internal/cskr/commands/splitCmd"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "cskr",
		Usage: "Curvy Social Key Recovery",
		Commands: []*cli.Command{
			splitCmd.NewCommand(),
			recoverCmd.NewCommand(),
		},
		ErrWriter: io.Discard,
	}

	logger := log.New(os.Stderr, "Error: ", 0)

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Fatal(err)
	}
}
