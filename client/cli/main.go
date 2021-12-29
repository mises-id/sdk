package main

import (
	cmd "github.com/mises-id/sdk/client/cli/commands"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/keys"
)

func main() {
	misesid.SetConfig()
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.LightCmd(),
		cmd.RestCmd(),
		keys.Commands(types.DefaultNodeHome),
		cli.NewCompletionCmd(rootCmd, true),
	)

	cmd := cli.PrepareBaseCmd(rootCmd, "Mises", types.DefaultNodeHome)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
