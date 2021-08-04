package user

import (
	//"encoding/json"
	"encoding/json"
	"fmt"

	//"io/ioutil"
	//"net/http"

	//"strings"
	"testing"

	"github.com/mises-id/sdk/bip39"
	"github.com/mises-id/sdk/types"
	"github.com/tyler-smith/assert"
)

/*
func TestCreateUser(t *testing.T) {
	entropy, err := bip39.NewEntropy(128)
	assert.NoError(t, err)

	mnemonics, err := bip39.NewMnemonic(entropy)
	assert.NoError(t, err)

	var ugr MisesUserMgr
	pUgr := &ugr
	user, err := pUgr.CreateUser(mnemonics, "123456")
	assert.NoError(t, err)

	session, err := CreateUser(user)
	fmt.Printf("session is: %s\n", session)
	assert.NoError(t, err)
}
*/
func CreateUserTest() types.MUser {
	//create user
	entropy, _ := bip39.NewEntropy(128)

	mnemonics, _ := bip39.NewMnemonic(entropy)

	var ugr MisesUserMgr
	pUgr := &ugr
	cuser, _ := pUgr.CreateUser(mnemonics, "123456")

	return cuser
}

/*
func TestWaitResp(t *testing.T) {
	//create user
	entropy, err := bip39.NewEntropy(128)
	assert.NoError(t, err)

	mnemonics, err := bip39.NewMnemonic(entropy)
	assert.NoError(t, err)

	var ugr MisesUserMgr
	pUgr := &ugr
	cuser, err := pUgr.CreateUser(mnemonics, "123456")
	assert.NoError(t, err)

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
	assert.NoError(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	fmt.Printf("%s task has been %s\n", r.session, string(body))
}

func TestGetUInfo(t *testing.T) {
	cuser := CreateUserTest()

	url, err := MakeGetUrl(cuser.MisesID(), GetUInfoUrl, cuser)
	assert.NoError(t, err)

	resp, err := http.Get(url)
	assert.NoError(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	fmt.Printf("body is %s\n", string(body))

	var uinfo MisesUserInfo
	err = json.Unmarshal(body, &uinfo)
	assert.NoError(t, err)

	ub, err := json.Marshal(&uinfo)
	assert.NoError(t, err)

	fmt.Printf("user info is %s\n", string(ub))
}


func TestGetFollowing(t *testing.T) {
	cuser := CreateUserTest()

	uf := cuser.MisesID()
	url, err := MakeGetUrl(uf, GetFollowingUrl, cuser)
	assert.NoError(t, err)

	resp, err := http.Get(url)
	assert.NoError(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	bs := string(body)
	fmt.Printf("body is %s\n", bs)

	followings := strings.Split(bs, "&")
	for _, f := range followings {
		fmt.Printf("%s ", f)
	}
	fmt.Printf("\n")
}


func TestSetFollowing(t *testing.T) {
	cuser := CreateUserTest()

	sessionid, err := SetFollowing(cuser, "followinguser", "follow")
	assert.NoError(t, err)

	fmt.Printf("following sessionid is %s\n", sessionid)

	time.Sleep(10000 * time.Millisecond)
}
*/

func TestCreateMisesID(t *testing.T) {
	cuser := CreateUserTest()

	sessionid, err := CreateUser(cuser)
	assert.NoError(t, err)
	assert.False(t, sessionid == "")
	fmt.Printf("create misesid sessionid is %s\n", sessionid)

}

func TestSetUserInfo(t *testing.T) {
	cuser := CreateUserTest()

	var info MisesUserInfo

	info.name = "yingming"
	info.gender = "男"
	info.avatarId = "007"
	info.avatarThumb = []byte("123456789")
	info.homePage = "http://mises.com"
	emails := []string{"yingming@gmail.com", "51911267@qq.com"}
	teles := []string{"17701314608", "18601350799", "18811790787"}
	info.emails = emails
	info.telephones = teles

	sessionid, err := SetUInfo(cuser, info)
	assert.NoError(t, err)
	assert.False(t, sessionid == "")

	fmt.Printf("userinfo sessionid is %s\n", sessionid)

}

func TestParseTxResp(t *testing.T) {

	resp := MsgTxResp{}
	resp.Code = 0
	resp.Error = ""
	resp.TxResponse = MsgTx{Height: "1", Txhash: "123456"}
	respBytes, err := json.Marshal(resp)
	fmt.Printf("resp is %s\n", string(respBytes))
	assert.NoError(t, err)
	msgTx, err := ParseTxResp(respBytes)
	assert.NoError(t, err)
	assert.EqualString(t, msgTx.Txhash, "123456")
}
