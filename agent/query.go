package sdk

/*
import (
	"encoding/json"
	"fmt"
)

//

type ReadResponseResult struct {
	Value string `json:"value"`
}

type ReadResponse struct {
	Result *ReadResponseResult `json:"result"`
}

func (ctx *misesAgent) Read(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf(ErrorKeyIsRequired)
	}
	if err := validateKey(key); err != nil {
		return "", err
	}

	body, err := ctx.APIQuery("/crud/read/" + ctx.uuid + "/" + encodeSafe(key))
	if err != nil {
		return "", err
	}

	res := &ReadResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return "", err
	}
	return res.Result.Value, nil
}

func (ctx *misesAgent) ProvenRead(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf(ErrorKeyIsRequired)
	}
	if err := validateKey(key); err != nil {
		return "", err
	}

	body, err := ctx.APIQuery("/crud/pread/" + ctx.uuid + "/" + encodeSafe(key))
	if err != nil {
		return "", err
	}

	res := &ReadResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return "", err
	}
	return res.Result.Value, nil
}

//

type Account struct {
	AccountNumber int     `json:"account_number"`
	Address       string  `json:"address"`
	Coins         []*Coin `json:"coins"`
	PublicKey     string  `json:"public_key"`
	Sequence      int     `json:"sequence"`
}

type Coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type AccountResponseResult struct {
	Value *Account `json:"value"`
}

type AccountResponse struct {
	Result *AccountResponseResult `json:"result"`
}

func (ctx *misesAgent) Account() (*Account, error) {
	res := &AccountResponse{}

	body, err := ctx.APIQuery("/auth/accounts/" + ctx.address)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res.Result.Value, nil
}

//

type VersionResponse struct {
	ApplicationVersion *VersionResponseApplicationVersion `json:"application_version"`
}

type VersionResponseApplicationVersion struct {
	Version string `json:"version"`
}

func (ctx *misesAgent) Version() (string, error) {
	body, err := ctx.APIQuery("/node_info")
	if err != nil {
		return "", err
	}

	res := &VersionResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return "", err
	}
	return res.ApplicationVersion.Version, nil
}
*/
