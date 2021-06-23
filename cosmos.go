package sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	tmcrypto "github.com/tendermint/tendermint/crypto"
	tmsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
)

const TX_COMMAND = "/txs"
const TOKEN_NAME = "ubnt"
const BROADCAST_MAX_RETRIES = 10
const BROADCAST_RETRY_INTERVAL = time.Second
const BLOCK_TIME_IN_SECONDS = 5

//
// JSON struct keys are ordered alphabetically
//

type ErrorResponse struct {
	Error string `json:"error"`
}

type KeyValue struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type KeyLease struct {
	Key   string `json:"key,omitempty"`
	Lease string `json:"lease,omitempty"`
}

type GasInfo struct {
	MaxGas   int `json:"max_gas"`
	MaxFee   int `json:"max_fee"`
	GasPrice int `json:"gas_price"`
}

type LeaseInfo struct {
	Days    int64 `json:"days"`
	Hours   int64 `json:"hours"`
	Minutes int64 `json:"minutes"`
	Seconds int64 `json:"seconds"`
}

func (lease *LeaseInfo) ToBlocks() int64 {
	var seconds int64
	seconds += lease.Days * 24 * 60 * 60
	seconds += lease.Hours * 60 * 60
	seconds += lease.Minutes * 60
	seconds += lease.Seconds
	return seconds / BLOCK_TIME_IN_SECONDS
}

//

type TransactionFeeAmount struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type TransactionFee struct {
	Amount []*TransactionFeeAmount `json:"amount"`
	Gas    string                  `json:"gas"`
}

//

type TransactionSignaturePubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type TransactionSignature struct {
	AccountNumber string                      `json:"account_number"`
	PubKey        *TransactionSignaturePubKey `json:"pub_key"`
	Sequence      string                      `json:"sequence"`
	Signature     string                      `json:"signature"`
}

//

type Transaction struct {
	Key       string
	KeyValues []*KeyValue
	Lease     int64
	N         uint64
	NewKey    string
	Value     string

	ApiRequestMethod   string
	ApiRequestEndpoint string
	GasInfo            *GasInfo

	done             chan bool
	result           []byte
	err              error
	broadcastRetries int
}

//

type TransactionValidateRequest struct {
	BaseReq   *TransactionValidateRequestBaseReq `json:"BaseReq"`
	Key       string                             `json:"Key,omitempty"`
	KeyValues []*KeyValue                        `json:"KeyValues,omitempty"`
	Lease     string                             `json:"Lease,omitempty"`
	N         string                             `json:"N,omitempty"`
	NewKey    string                             `json:"NewKey,omitempty"`
	Owner     string                             `json:"Owner"`
	UUID      string                             `json:"UUID"`
	Value     string                             `json:"Value,omitempty"`
}

type TransactionValidateRequestBaseReq struct {
	From    string `json:"from"`
	ChainId string `json:"chain_id"`
}

type TransactionValidateResponse struct {
	Type  string                       `json:"type"`
	Value *TransactionBroadcastPayload `json:"value"`
}

//

type TransactionBroadcastRequest struct {
	Transaction *TransactionBroadcastPayload `json:"tx"`
	Mode        string                       `json:"mode"`
}

type TransactionBroadcastResponse struct {
	Height    string `json:"height"`
	TxHash    string `json:"txhash"`
	Data      string `json:"data"`
	Codespace string `json:"codespace"`
	Code      int    `json:"code"`
	RawLog    string `json:"raw_log"`
	GasWanted string `json:"gas_wanted"`
}

//

type TransactionMsgValue struct {
	Key       string      `json:"Key,omitempty"`
	KeyValues []*KeyValue `json:"KeyValues,omitempty"`
	Lease     string      `json:"Lease,omitempty"`
	N         string      `json:"N,omitempty"`
	NewKey    string      `json:"NewKey,omitempty"`
	Owner     string      `json:"Owner"`
	UUID      string      `json:"UUID"`
	Value     string      `json:"Value,omitempty"`
}

type TransactionMsg struct {
	Type  string               `json:"type"`
	Value *TransactionMsgValue `json:"value"`
}

//

type TransactionBroadcastPayload struct {
	Fee        *TransactionFee         `json:"fee"`
	Memo       string                  `json:"memo"`
	Msg        []*TransactionMsg       `json:"msg"`
	Signatures []*TransactionSignature `json:"signatures"`
}

type TransactionBroadcastPayloadSignPayload struct {
	AccountNumber string            `json:"account_number"`
	ChainId       string            `json:"chain_id"`
	Fee           *TransactionFee   `json:"fee"`
	Memo          string            `json:"memo"`
	Msgs          []*TransactionMsg `json:"msgs"`
	Sequence      string            `json:"sequence"`
}

//

func (ctx *MisesAgent) APIQuery(endpoint string) ([]byte, error) {
	url := ctx.endpoint + endpoint

	ctx.Infof("get %s", url)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := parseResponse(res)
	return body, err
}

func (ctx *MisesAgent) APIMutate(method string, endpoint string, payload []byte) ([]byte, error) {
	url := ctx.endpoint + endpoint

	ctx.Infof("post %s", url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := parseResponse(res)
	return body, err
}

func (ctx *MisesAgent) SendTransaction(txn *Transaction) ([]byte, error) {
	txn.done = make(chan bool, 1)
	ctx.transactions <- txn
	done := <-txn.done
	if !done {
		ctx.Fatalf("txn did not complete") // todo: enqueue
	}
	if txn.err != nil {
		ctx.Errorf("transaction err(%s)", txn.err)
	}
	return txn.result, txn.err
}

func (ctx *MisesAgent) ProcessTransaction(txn *Transaction) {
	txn.broadcastRetries = 0

	var result []byte
	payload, err := ctx.ValidateTransaction(txn)
	if err == nil {
		result, err = ctx.BroadcastTransaction(payload, txn.GasInfo)
	}

	txn.result = result
	txn.err = err
	txn.done <- true
	close(txn.done)
}

// Get required min gas
func (ctx *MisesAgent) ValidateTransaction(txn *Transaction) (*TransactionBroadcastPayload, error) {
	req := &TransactionValidateRequest{
		BaseReq: &TransactionValidateRequestBaseReq{
			From:    ctx.address,
			ChainId: ctx.chainId,
		},
		UUID:      ctx.uuid,
		Key:       txn.Key,
		KeyValues: txn.KeyValues,
		Lease:     strconv.FormatInt(txn.Lease, 10),
		N:         strconv.FormatUint(txn.N, 10),
		NewKey:    txn.NewKey,
		Owner:     ctx.address,
		Value:     txn.Value,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx.Infof("txn init %+v", string(reqBytes))
	body, err := ctx.APIMutate(txn.ApiRequestMethod, txn.ApiRequestEndpoint, reqBytes)
	if err != nil {
		return nil, err
	}

	ctx.Infof("txn init %+v", string(body))

	res := &TransactionValidateResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

func (ctx *MisesAgent) BroadcastTransaction(txn *TransactionBroadcastPayload, gasInfo *GasInfo) ([]byte, error) {
	// Set memo
	txn.Memo = makeRandomString(32)

	// Set fee
	if gasInfo == nil {
		return nil, fmt.Errorf("gas_info is required")
	}
	gas, err := strconv.Atoi(txn.Fee.Gas)
	if err != nil {
		ctx.Errorf("failed to pass gas to int(%s)", txn.Fee.Gas)
	}
	amount := 0
	if len(txn.Fee.Amount) != 0 {
		if a, err := strconv.Atoi(txn.Fee.Amount[0].Amount); err == nil {
			amount = a
		}
	}
	if gasInfo.MaxGas != 0 && gas > gasInfo.MaxGas {
		gas = gasInfo.MaxGas
	}
	if gasInfo.MaxFee != 0 {
		amount = gasInfo.MaxFee
	} else if gasInfo.GasPrice != 0 {
		amount = gas * gasInfo.GasPrice
	}

	txn.Fee = &TransactionFee{
		Gas: strconv.Itoa(gas),
		Amount: []*TransactionFeeAmount{
			&TransactionFeeAmount{Denom: TOKEN_NAME, Amount: strconv.Itoa(amount)},
		},
	}

	// Set signatures
	if signature, err := ctx.SignTransaction(txn); err != nil {
		return nil, err
	} else {
		txn.Signatures = []*TransactionSignature{
			&TransactionSignature{
				PubKey: &TransactionSignaturePubKey{
					Type:  tmsecp256k1.PubKeyName,
					Value: base64.StdEncoding.EncodeToString(ctx.privateKey.PubKey().Bytes()),
				},
				Signature:     signature,
				AccountNumber: strconv.Itoa(ctx.account.AccountNumber),
				Sequence:      strconv.Itoa(ctx.account.Sequence),
			},
		}
	}

	// Broadcast txn
	req := &TransactionBroadcastRequest{
		Transaction: txn,
		Mode:        "block",
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx.Infof("txn broadcast request %+v", string(reqBytes))
	body, err := ctx.APIMutate("POST", TX_COMMAND, reqBytes)
	if err != nil {
		return nil, err
	}
	// ctx.Infof("txn broadcast response %+v", string(body))
	// Read txn broadcast response
	res := &TransactionBroadcastResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	ctx.Infof("txn broadcast response %+v", res)

	if res.Code == 0 {
		ctx.account.Sequence += 1
		if res.Data == "" {
			return []byte{}, nil
		}
		decodedData, err := hex.DecodeString(res.Data)
		return decodedData, err
	}
	if strings.Contains(res.RawLog, "signature verification failed") {
		ctx.broadcastRetries += 1
		ctx.Warnf("txn failed ... retrying(%d) ...", ctx.broadcastRetries)
		if ctx.broadcastRetries >= BROADCAST_MAX_RETRIES {
			return nil, fmt.Errorf("txn failed after max retry attempts")
		}
		time.Sleep(BROADCAST_RETRY_INTERVAL)
		// Lookup changed sequence
		if err := ctx.setAccount(); err != nil {
			return nil, err
		}
		b, err := ctx.BroadcastTransaction(txn, gasInfo)
		return b, err
	}

	return nil, fmt.Errorf("%s", res.RawLog)
}

func (ctx *MisesAgent) SignTransaction(txn *TransactionBroadcastPayload) (string, error) {
	payload := &TransactionBroadcastPayloadSignPayload{
		AccountNumber: strconv.Itoa(ctx.account.AccountNumber),
		ChainId:       ctx.chainId,
		Sequence:      strconv.Itoa(ctx.account.Sequence),
		Memo:          txn.Memo,
		Fee:           txn.Fee,
		Msgs:          txn.Msg,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	sanitized := sanitizeString(string(payloadBytes))
	ctx.Infof("txn sign %+v", sanitized)
	hash := tmcrypto.Sha256([]byte(sanitized))
	if s, err := ctx.privateKey.Sign(hash); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(s), nil
	}
}

func parseResponse(res *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	errRes := &ErrorResponse{}
	err = json.Unmarshal(body, errRes)
	if err != nil {
		return nil, err
	}

	if errRes.Error != "" {
		return nil, fmt.Errorf("%s", errRes.Error)
	}

	return body, nil
}
