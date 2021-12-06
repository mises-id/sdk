package mobile

import (
	"context"
	"fmt"
	"os"

	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"

	tmcfg "github.com/tendermint/tendermint/config"

	misescmd "github.com/mises-id/sdk/cmd/commands"
)

var _ MLightNode = &mLCD{}

var (
	NodeHome        string
	DefaultNodeHome string
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".misestm")
}

type mLCD struct {
}

func (lcd *mLCD) SetEndpoint(endpoint string) error {
	return nil
}

func (lcd *mLCD) Serve() error {
	_, err := CreateDefaultTendermintConfig(NodeHome)
	if err != nil {
		return err
	}

	interfaceRegistry := types.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(codec, tx.DefaultSignModes)
	clientCtx := client.Context{}.
		WithCodec(codec).
		WithHomeDir(NodeHome).
		WithTxConfig(txCfg).
		WithAccountRetriever(auth.AccountRetriever{})

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)
	cmd := misescmd.LightCmd()
	cmd.SetArgs([]string{
		"test",
		"--listening-address=tcp://0.0.0.0:26657",
		"--log-level=trace",
		"--primary-addr=http://e1.mises.site:26657",
		"--witness-addr=http://e2.mises.site:26657",
		"--trusted-height=582507",
		"--trusted-hash=3F541BDF3CF2CE414FB4A3FAF90931101C4ABD31093239AC7E7A787B3E387230",
		"--dir=" + NodeHome + "-light",
	})

	err = cmd.ExecuteContext(ctx)
	return err
}

func (lcd *mLCD) SetLogLevel(level int) error {
	return nil
}

func CreateDefaultTendermintConfig(rootDir string) (*tmcfg.Config, error) {
	conf := tmcfg.DefaultConfig()
	conf.SetRoot(rootDir)
	tmcfg.EnsureRoot(rootDir)

	if err := conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}

	return conf, nil
}

func NewMLightNode() MLightNode {
	return &mLCD{}
}
