package sdk

import (
	"fmt"
)

func (ctx *MisesAgent) Create(key string, value string, gasInfo *GasInfo, leaseInfo *LeaseInfo) error {
	if key == "" {
		return fmt.Errorf(ErrorKeyIsRequired)
	}
	if err := validateKey(key); err != nil {
		return err
	}
	if value == "" {
		return fmt.Errorf(ErrorValueIsRequired)
	}
	var lease int64
	if leaseInfo != nil {
		lease = leaseInfo.ToBlocks()
	}
	if lease < 0 {
		return fmt.Errorf(ErrorInvalidLeaseTime)
	}

	transaction := &Transaction{
		Key:                key,
		Value:              value,
		Lease:              lease,
		ApiRequestMethod:   "POST",
		ApiRequestEndpoint: "/crud/create",
		GasInfo:            gasInfo,
	}

	_, err := ctx.SendTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}
