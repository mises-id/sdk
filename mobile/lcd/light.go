package lcd

import (
	"context"
	"fmt"
	"os"

	tmcfg "github.com/tendermint/tendermint/config"

	"github.com/mises-id/sdk/client/cli/commands"
	"github.com/mises-id/sdk/types"
)

var _ MLightNode = &mLCD{}

type mLCD struct {
	chainId          string
	primaryAddress   string
	witnessAddresses string
	trustHeight      string
	trustHash        string
}

func (lcd *mLCD) SetChainID(chainId string) error {
	lcd.chainId = chainId
	return nil
}
func (lcd *mLCD) SetEndpoints(primary string, witnesses string) error {
	lcd.primaryAddress = primary
	lcd.witnessAddresses = witnesses
	return nil
}
func (lcd *mLCD) SetTrust(height string, hash string) error {
	lcd.trustHeight = height
	lcd.trustHash = hash
	return nil
}

// func (lcd *mLCD) ServeRestApi(listen string) error {
// 	_, err := CreateDefaultTendermintConfig(types.NodeHome)
// 	if err != nil {
// 		return err
// 	}

// 	interfaceRegistry := costypes.NewInterfaceRegistry()
// 	codec := codec.NewProtoCodec(interfaceRegistry)
// 	txCfg := tx.NewTxConfig(codec, tx.DefaultSignModes)
// 	clientCtx := client.Context{}.
// 		WithCodec(codec).
// 		WithHomeDir(types.NodeHome).
// 		WithTxConfig(txCfg).
// 		WithAccountRetriever(auth.AccountRetriever{}).
// 		WithInput(sdkrest.KeyringPass).
// 		WithKeyringDir("keyring")

// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)
// 	cmd := misescmd.RestCmd()
// 	cmd.SetArgs([]string{
// 		"--chain-id=" + lcd.chainId,
// 		"--listening-address=" + listen,
// 		"--log-level=trace",
// 	})

// 	err = cmd.ExecuteContext(ctx)
// 	return err
// }

func (lcd *mLCD) Serve(listen string) error {
	_, err := CreateDefaultTendermintConfig(types.NodeHome)
	if err != nil {
		return err
	}

	// interfaceRegistry := costypes.NewInterfaceRegistry()
	// codec := codec.NewProtoCodec(interfaceRegistry)
	// txCfg := tx.NewTxConfig(codec, tx.DefaultSignModes)
	// clientCtx := client.Context{}.
	// 	WithCodec(codec).
	// 	WithHomeDir(types.NodeHome).
	// 	WithTxConfig(txCfg).
	// 	WithAccountRetriever(auth.AccountRetriever{})

	ctx := context.Background()
	// ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)
	cmd := commands.LightCmd()
	cmd.SetArgs([]string{
		lcd.chainId,
		"--listening-address=" + listen,
		"--log-level=trace",
		"--primary-addr=" + lcd.primaryAddress,   //http://e1.mises.site:26657
		"--witness-addr=" + lcd.witnessAddresses, //http://e2.mises.site:26657
		"--trusted-height=" + lcd.trustHeight,    //963312
		"--trusted-hash=" + lcd.trustHash,        //219B062359064E5A00062624062D775C63AFFEC96361B40894F3C7B81437A660
		"--dir=" + types.NodeHome + "/light",
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

func SetHomePath(dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	types.NodeHome = dir + ".misestm"
	return nil
}
