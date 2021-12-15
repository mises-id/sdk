package main

import (
	cmd "github.com/mises-id/sdk/client/cli/commands"
	"github.com/mises-id/sdk/types"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/keys"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	AccountAddressPrefix = "mises"
)

var (
	AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"
)

func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}

func main() {
	SetConfig()
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
