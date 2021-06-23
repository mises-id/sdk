package sdk

import (
	tmsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
)

type misesAgent struct {
	MAgent
	address          string
	account          *Account
	logger           log.Logger
	privateKey       *tmsecp256k1.PrivKey
	broadcastRetries int
	transactions     chan *Transaction
	chainID          string
	endpoint         string
	uuid             string
	Debug            bool
}

func (ctx *misesAgent) setupLogger() {
	ctx.logger = log.NewNopLogger()
}

// Fetch the address account info (`number` and `sequence` to be used later)
func (ctx *misesAgent) setAccount() error {
	if account, err := ctx.Account(); err != nil {
		return err
	} else {
		ctx.account = account
		return nil
	}
}

func (ctx *misesAgent) processTransactions() {
	for txn := range ctx.transactions {
		// ctx.Infof("processing transaction(%+v)", txn)
		ctx.processTransaction(txn)
	}
}

func newMisesAgent(chainID string, seed *misesKeySeed) (*misesAgent, error) {

	ctx := &misesAgent{
		chainID: chainID,
	}

	ctx.setupLogger()
	pkey, err := seed.getPrivateKey()
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
