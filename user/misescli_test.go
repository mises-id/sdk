package user

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/mises-id/sdk/bip39"
	"github.com/tyler-smith/assert"
)

func TestCreateUser(t *testing.T) {
	entropy, err := bip39.NewEntropy(128)
	assert.Nil(t, err)

	mnemonics, err := bip39.NewMnemonic(entropy)
	assert.Nil(t, err)

	var ugr MisesUserMgr
	pUgr := &ugr
	user, err := pUgr.CreateUser(mnemonics, "123456")
	assert.Nil(t, err)

	session, err := CreateUser(user)
	fmt.Printf("session is: %s\n", session)
	assert.Nil(t, err)
}

func TestWaitResp(t *testing.T) {
	//create user
	entropy, err := bip39.NewEntropy(128)
	assert.Nil(t, err)

	mnemonics, err := bip39.NewMnemonic(entropy)
	assert.Nil(t, err)

	var ugr MisesUserMgr
	pUgr := &ugr
	cuser, err := pUgr.CreateUser(mnemonics, "123456")
	assert.Nil(t, err)

	// check mises decentralized result
	var r WaitResult
	r.session = "sessionid001"

	url, err := MakeGetUrl(r.session, QueryResultUrl, cuser)
	if err != nil {
		r.result = "601 url error"
		wr[r.session] = r
		return
	}

	resp, err := http.Get(url)
	assert.Nil(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	fmt.Printf("%s task has been %s\n", r.session, string(body))
}
