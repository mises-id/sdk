package main

import (
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	cmd "github.com/mises-id/sdk/client/cli/commands"
	light "github.com/mises-id/sdk/client/cli/commands/light"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
	"github.com/tendermint/tendermint/libs/cli"
)

func main() {
	misesid.SetConfig()
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		light.LightCmd(),
		cmd.KeysCmd(types.DefaultNodeHome),
		cli.NewCompletionCmd(rootCmd, true),
	)

	//cmd := cli.PrepareBaseCmd(rootCmd, "Mises", types.DefaultNodeHome)
	if err := svrcmd.Execute(rootCmd, types.DefaultNodeHome); err != nil {
		panic(err)
	}
}
