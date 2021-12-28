/*
	App is a Community Site which supports misesid,
	there is only one app "mises.site" during the first stage(MVP)
*/
package app

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcodec "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdkrest "github.com/mises-id/sdk/client/rest"

	misestypes "github.com/mises-id/mises-tm/x/misestm/types"
	cmd "github.com/mises-id/sdk/client/cli/commands"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
)

type MisesApp struct {
	info   types.MAppInfo
	appDid string
	auths  []MisesAuth

	clientCtx client.Context

	pubKey string
}

type MisesAuth struct {
	Uid                 string
	ExpirationInSeconds int
	Permissions         []string
}

const (
	MisesDiscover              = "mises.site"
	MisesDiscoverAppKey        = "mises-discover"
	defaultExpirationInSeconds = 120
)

func (app *MisesApp) AppDID() string {
	return app.appDid
}
func (app *MisesApp) MisesID() string {
	return app.appDid
}

func (app *MisesApp) Info() types.MAppInfo {
	return app.info
}

func (app *MisesApp) AddAuth(misesId string, permissions []string) {
	var auth MisesAuth
	auth.Uid = misesId
	auth.ExpirationInSeconds = defaultExpirationInSeconds // default is 120 seconds
	auth.Permissions = permissions

	app.auths = append(app.auths, auth)
}

func (app *MisesApp) Init(chainID string, passPhrase string) error {
	cmd.SetConfig()
	clientCtx := client.Context{}
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	authcodec.RegisterInterfaces(interfaceRegistry)
	bankcodec.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)

	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(codec, tx.DefaultSignModes)

	clientCtx = clientCtx.
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithCodec(codec).
		WithInterfaceRegistry(interfaceRegistry).
		WithTxConfig(txCfg).
		WithInput(&sdkrest.PassReader{Pass: passPhrase}).
		WithKeyringDir(types.NodeHome + "/sdk-keyring").
		WithChainID(chainID)
	kr, err := client.NewKeyringFromBackend(clientCtx, keyring.BackendFile)
	if err != nil {
		return err
	}

	key, err := kr.Key(MisesDiscoverAppKey)
	if err != nil {
		mnemonics, err := misesid.RandomMnemonics()
		if err != nil {
			return err
		}
		fmt.Printf("app mnemonics is: %s\n", mnemonics)
		key, err = kr.NewAccount(MisesDiscoverAppKey, mnemonics, passPhrase, "", hd.Secp256k1)
		if err != nil {
			return err
		}

	}

	rpcURI := "tcp://127.0.0.1:26657"

	clientCtx = clientCtx.WithNodeURI(rpcURI)

	client, err := client.NewClientFromNode(rpcURI)
	if err != nil {
		return err
	}

	clientCtx = clientCtx.WithClient(client)

	clientCtx = clientCtx.WithFromName(key.GetName()).WithFromAddress(key.GetAddress())
	clientCtx = clientCtx.WithKeyring(kr)

	if err := clientCtx.AccountRetriever.EnsureExists(clientCtx, key.GetAddress()); err != nil {
		return err
	}

	app.clientCtx = clientCtx

	app.pubKey = hex.EncodeToString(key.GetPubKey().Bytes())

	app.appDid = types.MisesAppIDPrefix + key.GetAddress().String()
	app.info = misesid.NewMisesAppInfoReadonly(
		"Mises Discover'",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{MisesDiscover},
		"Mises Network",
	)

	if err := misesid.StarSeqGenerator(app.clientCtx); err != nil {
		return err
	}

	if _, err := misesid.GetMisesID(app, app.MisesID()); err != nil {

		tx, err := misesid.CreateDid(app.clientCtx, app.pubKey, app.MisesID())
		if err != nil {
			return err
		}
		err = misesid.PollTxSync(app.clientCtx, tx)
		if err != nil {
			return err
		}
		tx, err = misesid.UpdateAppInfo(app.clientCtx, app.MisesID(), misestypes.PublicAppInfo{
			Name:      app.info.AppName(),
			Domains:   app.info.Domains(),
			Developer: app.info.Developer(),
			HomeUrl:   app.info.HomeURL(),
			IconUrl:   app.info.IconURL(),
		})
		if err != nil {
			return err
		}
		err = misesid.PollTxSync(app.clientCtx, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *MisesApp) RegisterUser(misesUID string, userPubKey string) error {
	if _, err := misesid.GetMisesID(app, misesUID); err != nil {

		tx, err := misesid.CreateDid(app.clientCtx, userPubKey, misesUID)
		if err != nil {
			return err
		}
		err = misesid.PollTxSync(app.clientCtx, tx)
		if err != nil {
			return err
		}
	}
	tx, err := misesid.UpdateAppFeeGrant(app.clientCtx, app.MisesID(), misesUID)
	if err != nil {
		return err
	}
	err = misesid.PollTxSync(app.clientCtx, tx)
	if err != nil {
		return err
	}
	return nil

}

func (app *MisesApp) PubKey() string {
	return app.pubKey
}

func (app *MisesApp) Sign(msg string) (string, error) {

	kr := app.clientCtx.Keyring
	if kr == nil {
		return "", fmt.Errorf("no keyring")
	}

	sigBytes, _, err := kr.Sign(MisesDiscoverAppKey, []byte(msg))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sigBytes), nil

}

func (app *MisesApp) AesKey() ([]byte, error) {
	kr := app.clientCtx.Keyring
	if kr == nil {
		return nil, fmt.Errorf("no keyring")
	}
	privKey, err := keyring.NewUnsafe(kr).UnsafeExportPrivKeyHex(MisesDiscoverAppKey)
	if err != nil {
		return nil, err
	}
	privKeyBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, err
	}
	mhash := sha256.Sum256([]byte(privKeyBytes))
	return mhash[:], nil
}

func (app *MisesApp) Signer() types.MSigner {
	return app
}
