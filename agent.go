package sdk

import (
	tmsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
)

const DEFAULT_ENDPOINT string = "http://localhost:1317"
const DEFAULT_CHAIN_ID string = "mises"
const HD_PATH string = "m/44'/118'/0'/0/0"
const ADDRESS_PREFIX string = "mises"

type MisesAgent struct {
	MAgent
	address          string
	account          *Account
	logger           log.Logger
	privateKey       *tmsecp256k1.PrivKey
	broadcastRetries int
	transactions     chan *Transaction
	chainId          string
	endpoint         string
	uuid             string
	Debug            bool
}

func (ctx *MisesAgent) setupLogger() {
	ctx.logger = log.NewNopLogger()
}

// Fetch the address account info (`number` and `sequence` to be used later)
func (ctx *MisesAgent) setAccount() error {
	if account, err := ctx.Account(); err != nil {
		return err
	} else {
		ctx.account = account
		return nil
	}
}

func (ctx *MisesAgent) processTransactions() {
	for txn := range ctx.transactions {
		// ctx.Infof("processing transaction(%+v)", txn)
		ctx.ProcessTransaction(txn)
	}
}

func newMisesAgent(chain_id string, seed *MisesKeySeed) (*MisesAgent, error) {

	ctx := &MisesAgent{
		chainId: chain_id,
	}

	ctx.setupLogger()
	pkey, err := seed.genPrivateKey()
	// Generate private key from mnemonic
	if err != nil {
		return nil, err
	}
	ctx.privateKey = pkey

	addr, err := seed.getAddress()
	if err != nil {
		return nil, err
	}
	ctx.address = addr

	// Fetch the address account info (`number` and `sequence` to be used later)
	if err := ctx.setAccount(); err != nil {
		return nil, err
	}

	// Send transactions
	ctx.transactions = make(chan *Transaction, 1) // serial
	go ctx.processTransactions()

	return ctx, nil
}
