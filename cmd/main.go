package main

import (
	"os"
	"path/filepath"

	cmd "github.com/mises-id/sdk/cmd/commands"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
)

func main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.LightCmd(),
		cmd.RestCmd(),
		cli.NewCompletionCmd(rootCmd, true),
	)

	cmd := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("$HOME", cfg.DefaultTendermintDir)))
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
