package misesid

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/mises-id/mises-tm/x/misestm/types"
	multibase "github.com/multiformats/go-multibase"
)

type SeqInfo struct {
	nextNum uint64
	nextSeq uint64
}

var (
	SeqCmdChan  chan int
	SeqInfoChan chan SeqInfo
)

func PollTxSync(clientCtx client.Context, tx *sdk.TxResponse) (err error) {
	if tx.Code != 0 {
		return fmt.Errorf("tx fail: %s", tx.RawLog)
	}
	var errCount int = 0
	for {

		err = PollTx(clientCtx, tx.TxHash)

		if err != nil {

			if errCount > 10 {
				return err
			}
			errCount += 1
			time.Sleep(2 * time.Second)
			continue
		} else {
			break
		}

	}
	return
}
func PollTx(clientCtx client.Context, txHash string) error {
	hash, err := hex.DecodeString(txHash)
	if err != nil {
		return err
	}

	node, err := clientCtx.GetNode()
	if err != nil {
		return err
	}

	resTx, err := node.Tx(context.Background(), hash, true)
	if err != nil {
		return err
	}

	if resTx.Height == 0 {
		return fmt.Errorf("tx fail [" + txHash + "] " + resTx.TxResult.Log)
	}

	return nil
}

func StarSeqGenerator(clientCtx client.Context) error {

	seqChan := make(chan SeqInfo, 1)
	cmdChan := make(chan int, 1)
	ar := clientCtx.AccountRetriever
	var key keyring.Info
	var err error
	if clientCtx.Keyring != nil {
		key, err = clientCtx.Keyring.KeyByAddress(clientCtx.FromAddress)
	} else {
		return fmt.Errorf("no key ring")
	}
	if err != nil {
		return err
	}
	keyaddr := key.GetAddress()

	go func() {
		var num, seq uint64
		for {
			if seq == 0 {
				var err error
				num, seq, err = ar.GetAccountNumberSequence(clientCtx, keyaddr)
				if err != nil {
					time.Sleep(2 * time.Second)
					continue
				}
			}

			seqChan <- SeqInfo{nextNum: num, nextSeq: seq}

			next := <-cmdChan
			if next == 1 {
				seq = 0
			} else {
				seq++
			}

		}

	}()
	SeqInfoChan = seqChan
	SeqCmdChan = cmdChan
	return nil
}

func prepareFactory(clientCtx client.Context, txf tx.Factory) tx.Factory {
	gasSetting := flags.GasSetting{
		Simulate: false,
		Gas:      100000,
	}
	txf = txf.
		WithTxConfig(clientCtx.TxConfig).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithKeybase(clientCtx.Keyring).
		WithChainID(clientCtx.ChainID).
		WithGas(gasSetting.Gas).
		WithSimulateAndExecute(gasSetting.Simulate).
		WithTimeoutHeight(0).
		WithGasAdjustment(2.0).
		WithMemo("mises go sdk").
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	seq := <-SeqInfoChan
	txf = txf.WithAccountNumber(seq.nextNum)
	txf = txf.WithSequence(seq.nextSeq)

	return txf
}

// CalculateGas simulates the execution of a transaction and returns the
// simulation response obtained by the query and the adjusted gas amount.
func calculateGas(
	clientCtx client.Context, txf tx.Factory, msgs ...sdk.Msg,
) (*sdk.Result, uint64, error) {
	txBytes, err := txf.BuildSimTx(msgs...)
	if err != nil {
		return nil, 0, err
	}

	gasInfo, result, err := types.Simulater(txBytes)
	if err != nil {
		return nil, 0, err
	}

	return result, uint64(txf.GasAdjustment() * float64(gasInfo.GasUsed)), nil
}

func broadcastTx(clientCtx client.Context, txf tx.Factory, msgs ...sdk.Msg) (*sdk.TxResponse, error) {

	if txf.SimulateAndExecute() || clientCtx.Simulate {
		_, adjusted, err := calculateGas(clientCtx, txf, msgs...)
		if err != nil {
			return nil, err
		}

		txf = txf.WithGas(adjusted)
	}
	if clientCtx.Simulate {
		return nil, nil
	}

	txb, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	if clientCtx.GetFeeGranterAddress() != nil {
		txb.SetFeeGranter(clientCtx.GetFeeGranterAddress())
	}

	err = authclient.SignTx(txf, clientCtx, clientCtx.GetFromName(), txb, true, true)
	if err != nil {
		return nil, err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		return nil, err
	}

	//types.Logger.Error(fmt.Sprintf("BroadcastTx start with seq %v", txf.Sequence()))

	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, err
	}

	//types.Logger.Error(fmt.Sprintf("BroadcastTx finish with code %v", res.Code))

	return postBroadcastTx(clientCtx, res, err)
}

func prepareSigner(clientCtx client.Context) (client.Context, error) {
	if clientCtx.ChainID == "" {
		clientCtx = clientCtx.WithChainID("mises")
	}
	if clientCtx.Keyring == nil {
		panic(fmt.Errorf("no key ring"))
	}

	clientCtx = clientCtx.WithBroadcastMode(flags.BroadcastSync)
	return clientCtx, nil
}

func CreateDid(clientCtx client.Context, pubKeyHex string, misesID string) (*sdk.TxResponse, error) {
	clientCtx, err := prepareSigner(clientCtx)
	if err != nil {
		return nil, err
	}
	if err := CheckMisesID(misesID, pubKeyHex); err != nil {
		return nil, err
	}

	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, err
	}
	pubKeyMultiBase, err := multibase.Encode(multibase.Base58BTC, pubKeyBytes)
	if err != nil {
		return nil, err
	}

	msg := types.NewMsgCreateDidRegistry(
		clientCtx.FromAddress.String(),
		misesID,
		misesID+"#key0",
		"EcdsaSecp256k1VerificationKey2019", // will shift to Ed25519VerificationKey2020
		pubKeyMultiBase,
		0,
	)
	if err := msg.ValidateBasic(); err != nil {
		if err != nil {
			return nil, err
		}

	}
	txf := tx.Factory{}
	txf = prepareFactory(clientCtx, txf)

	return broadcastTx(clientCtx, txf, msg)
}

func postBroadcastTx(clientCtx client.Context, res *sdk.TxResponse, err error) (*sdk.TxResponse, error) {

	if err != nil {
		SeqCmdChan <- -1
		return nil, err
	}

	if res.Code == sdkerrors.ErrWrongSequence.ABCICode() {
		//reset cmdSeqChan
		SeqCmdChan <- 1
	} else {
		SeqCmdChan <- 0
	}

	return res, nil
}

func UpdateUserInfo(clientCtx client.Context, misesUid string, priInfo types.PrivateUserInfo) (*sdk.TxResponse, error) {

	clientCtx, err := prepareSigner(clientCtx)
	if err != nil {
		return nil, err
	}

	msg := types.NewMsgUpdateUserInfo(
		clientCtx.FromAddress.String(),
		misesUid,
		priInfo.EncData,
		priInfo.Iv,
		0,
	)
	if err := msg.ValidateBasic(); err != nil {
		if err != nil {
			return nil, err
		}

	}
	txf := tx.Factory{}
	txf = prepareFactory(clientCtx, txf)

	return broadcastTx(clientCtx, txf, msg)
}

func UpdateUserRelation(clientCtx client.Context, actionStr string, misesUID string, targetUID string) (*sdk.TxResponse, error) {

	clientCtx, err := prepareSigner(clientCtx)
	if err != nil {
		return nil, err
	}

	var action uint64
	switch actionStr {
	case "follow":
		action = 0
	case "unfollow":
		action = 1
	case "block":
		action = 2
	case "unblock":
		action = 3
	}

	msg := types.NewMsgUpdateUserRelation(
		clientCtx.FromAddress.String(),
		misesUID,
		targetUID,
		action,
		0,
	)
	if err := msg.ValidateBasic(); err != nil {
		if err != nil {
			return nil, err
		}

	}
	txf := tx.Factory{}
	txf = prepareFactory(clientCtx, txf)

	return broadcastTx(clientCtx, txf, msg)
}

func UpdateAppInfo(clientCtx client.Context, misesAppID string, pubInfo types.PublicAppInfo) (*sdk.TxResponse, error) {

	clientCtx, err := prepareSigner(clientCtx)
	if err != nil {
		return nil, err
	}

	msg := types.NewMsgUpdateAppInfo(
		clientCtx.FromAddress.String(),
		misesAppID,
		pubInfo.Name,
		pubInfo.Domains,
		pubInfo.Developer,
		pubInfo.HomeUrl,
		pubInfo.IconUrl,
		0,
	)
	if err := msg.ValidateBasic(); err != nil {
		if err != nil {
			return nil, err
		}
	}
	txf := tx.Factory{}
	txf = prepareFactory(clientCtx, txf)

	return broadcastTx(clientCtx, txf, msg)
}

func UpdateAppFeeGrant(clientCtx client.Context, misesAppID string, misesUid string) (*sdk.TxResponse, error) {

	clientCtx, err := prepareSigner(clientCtx)
	if err != nil {
		return nil, err
	}

	basic := feegrant.BasicAllowance{}

	periodLimit := []sdk.Coin{{
		Denom:  "umis",
		Amount: sdk.NewInt(10000),
	}}
	period := time.Duration(24 * 3600 * 1000000000) //1day

	periodic := feegrant.PeriodicAllowance{
		Basic:            basic,
		Period:           period,
		PeriodReset:      time.Now().Add(period),
		PeriodSpendLimit: periodLimit,
		PeriodCanSpend:   periodLimit,
	}

	var allowance feegrant.FeeAllowanceI
	allowance = &periodic

	allowedMsgs := []string{
		"/misesid.misestm.v1beta1.MsgUpdateUserInfo",
		"/misesid.misestm.v1beta1.MsgUpdateUserRelation",
	}

	allowance, err = feegrant.NewAllowedMsgAllowance(allowance, allowedMsgs)
	if err != nil {
		return nil, err
	}

	appAddr, _, err := types.AddrFromDid(misesAppID)
	if err != nil {
		return nil, err
	}

	userAddr, _, err := types.AddrFromDid(misesUid)
	if err != nil {
		return nil, err
	}

	msg, err := feegrant.NewMsgGrantAllowance(allowance, appAddr, userAddr)
	if err != nil {
		return nil, err
	}

	if err := msg.ValidateBasic(); err != nil {
		if err != nil {
			return nil, err
		}

	}
	txf := tx.Factory{}
	txf = prepareFactory(clientCtx, txf)

	return broadcastTx(clientCtx, txf, msg)
}
