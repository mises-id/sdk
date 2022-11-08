/*
	App is a Community Site which supports misesid,
	there is only one app "mises.site" during the first stage(MVP)
*/
package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcodec "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	misestypes "github.com/mises-id/mises-tm/x/misestm/types"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"

	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "app")
)

type MisesApp struct {
	info   types.MAppInfo
	appDid string
	auths  []MisesAuth

	clientCtx client.Context

	pubKey string

	pendingCmds chan types.MisesAppCmd

	waitingCmds chan types.MisesAppCmd

	failedTxCounter map[string]int

	listener types.MisesAppCmdListener

	seqChan *misesid.SeqChan
}

type MisesAppCmdBase struct {
	misesUID string
	pubKey   string
	txid     string
	waitTx   bool
	trackID  string
}

func (cmd *MisesAppCmdBase) MisesUID() string {
	return cmd.misesUID
}
func (cmd *MisesAppCmdBase) PubKey() string {
	return cmd.pubKey
}

func (cmd *MisesAppCmdBase) TxID() string {
	return cmd.txid
}

func (cmd *MisesAppCmdBase) SetTxID(txid string) {
	cmd.txid = txid
}

func (cmd *MisesAppCmdBase) TrackID() string {
	return cmd.trackID
}

func (cmd *MisesAppCmdBase) SetTrackID(trackID string) {
	cmd.trackID = trackID
}

func (cmd *MisesAppCmdBase) WaitTx() bool {
	return cmd.waitTx
}

func (cmd *MisesAppCmdBase) SetWaitTx(waitTx bool) {
	cmd.waitTx = waitTx
}

type RegisterUserCmd struct {
	MisesAppCmdBase
	feeGrantedPerDay int64
}

func (cmd *RegisterUserCmd) FeeGrantedPerDay() int64 {
	return cmd.feeGrantedPerDay
}

type FaucetCmd struct {
	MisesAppCmdBase
	coinUMIS int64
}

func (cmd *FaucetCmd) CoinUMIS() int64 {
	return cmd.coinUMIS
}

type MisesAuth struct {
	Uid                 string
	ExpirationInSeconds int
	Permissions         []string
}
type PassReader struct {
	Pass string
}

func (r *PassReader) Read(p []byte) (n int, err error) {
	n = copy(p, []byte(r.Pass))
	n += copy(p[n:], []byte("\n"))
	return
}

const (
	defaultExpirationInSeconds = 120
	maxPendingCmds             = 100
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

func (app *MisesApp) AppKey(name string) string {
	key := strings.ToLower(name)
	key = strings.ReplaceAll(key, " ", "-")
	return key
}

func (app *MisesApp) AddAuth(misesId string, permissions []string) {
	var auth MisesAuth
	auth.Uid = misesId
	auth.ExpirationInSeconds = defaultExpirationInSeconds // default is 120 seconds
	auth.Permissions = permissions

	app.auths = append(app.auths, auth)
}

func (app *MisesApp) Init(info types.MAppInfo, options types.MSdkOption) error {
	misesid.SetConfig()
	clientCtx := client.Context{}
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	authcodec.RegisterInterfaces(interfaceRegistry)
	bankcodec.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	misestypes.RegisterInterfaces(interfaceRegistry)

	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(codec, tx.DefaultSignModes)

	clientCtx = clientCtx.
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithCodec(codec).
		WithInterfaceRegistry(interfaceRegistry).
		WithTxConfig(txCfg).
		WithInput(&PassReader{Pass: options.PassPhrase}).
		WithKeyringDir(types.NodeHome + "/sdk-keyring").
		WithChainID(options.ChainID)
	kr, err := client.NewKeyringFromBackend(clientCtx, keyring.BackendFile)
	if err != nil {
		return err
	}

	appKey := app.AppKey(info.AppName())
	key, err := kr.Key(appKey)
	if err != nil {
		mnemonics, err := misesid.RandomMnemonics()
		if err != nil {
			return err
		}
		logger.Info("app mnemonics is: ", mnemonics)
		key, err = kr.NewAccount(appKey, mnemonics, options.PassPhrase, "", hd.Secp256k1)
		if err != nil {
			return err
		}
		logger.Info("app address is: ", key.GetAddress().String())

	}

	clientCtx = clientCtx.WithNodeURI(options.RpcURI)

	client, err := client.NewClientFromNode(options.RpcURI)
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
	app.info = info

	if app.seqChan, err = misesid.StarSeqGenerator(app.clientCtx, logger); err != nil {
		return err
	}

	if err := misesid.CheckDid(app.clientCtx, app.MisesID()); err != nil {

		tx, err := misesid.CreateDid(app.clientCtx, app.seqChan, app.pubKey, app.MisesID())
		if err != nil {
			return err
		}
		err = misesid.PollTxSync(app.clientCtx, tx)
		if err != nil {
			return err
		}
		tx, err = misesid.UpdateAppInfo(app.clientCtx, app.seqChan, app.MisesID(), misestypes.PublicAppInfo{
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

	app.startCmdRoutine()

	return nil
}

func (app *MisesApp) asynWaitCmd(cmd types.MisesAppCmd, err error) {
	if cmd.TxID() == "" {
		if app.listener != nil {
			app.listener.OnFailed(cmd, fmt.Errorf("no txid"))
		}
		return
	}
	if failCount, ok := app.failedTxCounter[cmd.TxID()]; ok {
		//do something here
		if failCount > 10 {
			delete(app.failedTxCounter, cmd.TxID())
			if app.listener != nil {
				if err != nil {
					app.listener.OnFailed(cmd, fmt.Errorf("fail after 10 try: %s", err.Error()))
				} else {
					app.listener.OnFailed(cmd, fmt.Errorf("fail after 10 try"))
				}

			}
			return
		}
		app.failedTxCounter[cmd.TxID()] = failCount + 1
	} else {
		app.failedTxCounter[cmd.TxID()] = 0
	}
	go func(cmd types.MisesAppCmd) {
		time.Sleep(2 * time.Second)
		app.waitingCmds <- cmd
	}(cmd)
}
func (app *MisesApp) startCmdRoutine() {
	app.pendingCmds = make(chan types.MisesAppCmd, maxPendingCmds)
	app.waitingCmds = make(chan types.MisesAppCmd, maxPendingCmds)
	app.failedTxCounter = map[string]int{}
	go func() {
		for {

			cmd := <-app.pendingCmds

			err := app.RunSync(cmd)
			if err != nil {
				logger.Info("cmd fail: ", err.Error())
				if app.listener != nil {
					app.listener.OnFailed(cmd, err)
				}
			} else {
				if cmd.WaitTx() {
					if app.listener != nil {
						app.listener.OnSucceed(cmd)
					}
				} else {
					app.asynWaitCmd(cmd, nil)
				}

			}

		}
	}()

	go func() {
		for {

			cmd := <-app.waitingCmds

			resTx, err := misesid.PollTx(app.clientCtx, cmd.TxID())
			if err != nil {
				app.asynWaitCmd(cmd, err)
			} else {
				delete(app.failedTxCounter, cmd.TxID())
				if resTx.Height == 0 || resTx.TxResult.Code != 0 {
					if app.listener != nil {
						app.listener.OnFailed(cmd, fmt.Errorf("fail with code %d", resTx.TxResult.Code))
					}
				} else {
					if app.listener != nil {
						app.listener.OnSucceed(cmd)
					}
				}
			}

		}
	}()
}

func (app *MisesApp) RunSync(cmd types.MisesAppCmd) error {
	//ensure did
	if err := misesid.CheckDid(app.clientCtx, cmd.MisesUID()); err != nil {
		if cmd.PubKey() == "" {
			return fmt.Errorf("no pubkey")
		}

		tx, err := misesid.CreateDid(app.clientCtx, app.seqChan, cmd.PubKey(), cmd.MisesUID())
		if err != nil {
			return err
		}
		if app.listener != nil {
			cmd.SetTxID(tx.TxHash)
			app.listener.OnTxGenerated(cmd)
		}
		if cmd.WaitTx() {
			err = misesid.PollTxSync(app.clientCtx, tx)
			if err != nil {
				return err
			}
		}

	}

	var tx *sdk.TxResponse = nil
	var err error
	if cmdapp, ok := cmd.(*RegisterUserCmd); ok {
		if err := misesid.CheckAppFeeGrant(app.clientCtx, app.MisesID(), cmdapp.MisesUID()); err != nil {
			tx, err = misesid.UpdateAppFeeGrant(app.clientCtx, app.seqChan, app.MisesID(), cmdapp.MisesUID(), cmdapp.FeeGrantedPerDay())
		}
	} else if cmdapp, ok := cmd.(*FaucetCmd); ok {
		tx, err = misesid.Transfer(app.clientCtx, app.seqChan, app.MisesID(), cmdapp.MisesUID(), cmdapp.CoinUMIS())
	} else {
		return fmt.Errorf("known cmd")
	}
	if err != nil {
		return err
	}

	if tx != nil {
		if app.listener != nil {
			cmd.SetTxID(tx.TxHash)
			app.listener.OnTxGenerated(cmd)
		}
		if cmd.WaitTx() {
			err = misesid.PollTxSync(app.clientCtx, tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *MisesApp) RunAsync(cmd types.MisesAppCmd, wait bool) error {
	if len(app.pendingCmds) == maxPendingCmds {
		return fmt.Errorf("too many pending commands")
	}
	cmd.SetWaitTx(wait)
	app.pendingCmds <- cmd

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
	appKey := app.AppKey(app.info.AppName())

	sigBytes, _, err := kr.Sign(appKey, []byte(msg))
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
	appKey := app.AppKey(app.info.AppName())
	privKey, err := keyring.NewUnsafe(kr).UnsafeExportPrivKeyHex(appKey)
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

func (app *MisesApp) NewRegisterUserCmd(uid string, pubkey string, feeGrantedPerDay int64) types.MisesAppCmd {
	return &RegisterUserCmd{
		MisesAppCmdBase{uid, pubkey, "", true, ""}, feeGrantedPerDay,
	}
}
func (app *MisesApp) NewFaucetCmd(uid string, pubkey string, coinUMIS int64) types.MisesAppCmd {
	return &FaucetCmd{
		MisesAppCmdBase{uid, pubkey, "", true, ""}, coinUMIS,
	}
}

func (app *MisesApp) SetListener(listener types.MisesAppCmdListener) {
	app.listener = listener
}

func (app *MisesApp) StartEventStreaming(listener types.MisesEventStreamingListener) error {
	clientCtx := app.clientCtx
	var err error
	if !clientCtx.Client.IsRunning() {
		err = clientCtx.Client.Start()
		if err != nil {
			return err
		}
	}

	err = clientCtx.Client.UnsubscribeAll(context.Background(), app.appDid)
	if err != nil {
		return err
	}

	newBlockChan, err := clientCtx.Client.Subscribe(
		context.Background(),
		app.appDid,
		"tm.event = 'NewBlockHeader'",
		1,
	)
	if err != nil {
		return err
	}

	go func() {
		for event := range newBlockChan {
			if newBlockHeaderEvent, ok := event.Data.(tmtypes.EventDataNewBlockHeader); ok {
				listener.OnNewBlockHeaderEvent(&newBlockHeaderEvent)
			}
		}

		listener.OnEventStreamingTerminated()

	}()

	txChan, err := clientCtx.Client.Subscribe(
		context.Background(),
		app.appDid,
		"tm.event = 'Tx'",
		1,
	)
	if err != nil {
		return err
	}

	go func() {
		for event := range txChan {
			if txEvent, ok := event.Data.(tmtypes.EventDataTx); ok {

				listener.OnTxEvent(&txEvent)
			}
		}

		listener.OnEventStreamingTerminated()

	}()

	return nil
}

type intoAny interface {
	AsAny() *codectypes.Any
}

func (app *MisesApp) ParseEvent(header *tmtypes.EventDataNewBlockHeader, tx *tmtypes.EventDataTx) (*sdk.TxResponse, error) {
	clientCtx := app.clientCtx
	txb, err := clientCtx.TxConfig.TxDecoder()(tx.Tx)
	if err != nil {
		return nil, err
	}
	p, ok := txb.(intoAny)
	if !ok {
		return nil, fmt.Errorf("expecting a type implementing intoAny, got: %T", txb)
	}
	anyTx := p.AsAny()

	parsedLogs, _ := sdk.ParseABCILogs(tx.Result.Log)

	return &sdk.TxResponse{
		TxHash:    bytes.HexBytes(tmtypes.Tx(tx.Tx).Hash()).String(),
		Height:    header.Header.Height,
		Codespace: tx.Result.Codespace,
		Code:      tx.Result.Code,
		Data:      strings.ToUpper(hex.EncodeToString(tx.Result.Data)),
		RawLog:    tx.Result.Log,
		Logs:      parsedLogs,
		Info:      tx.Result.Info,
		GasWanted: tx.Result.GasWanted,
		GasUsed:   tx.Result.GasUsed,
		Tx:        anyTx,
		Timestamp: header.Header.Time.Format(time.RFC3339),
		Events:    tx.Result.Events,
	}, nil
}
