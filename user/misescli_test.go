package user_test

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"time"

	//"io/ioutil"
	//"net/http"

	//"strings"
	"testing"

	"github.com/mises-id/sdk/bip39"
	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sdk/user"
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

	var ugr user.MisesUserMgr
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

func PollSession(t *testing.T, session string) {
	wr, err := user.PollSessionResult(60 * time.Second)
	fmt.Printf("PollSessionResult finish\n")
	assert.NoError(t, err)
	assert.True(t, wr.ErrMsg == "")
	assert.True(t, wr.Session == session)
}
func PrepareUser(t *testing.T, cuser types.MUser) {
	sessionid, err := user.CreateUser(cuser)
	assert.NoError(t, err)
	assert.False(t, sessionid == "")
	fmt.Printf("create misesid sessionid is %s\n", sessionid)
	PollSession(t, sessionid)
}

func TestCreateMisesID(t *testing.T) {
	cuser := CreateUserTest()

	PrepareUser(t, cuser)

}

func TestSetUserInfo(t *testing.T) {
	cuser := CreateUserTest()

	PrepareUser(t, cuser)

	info := user.NewMisesUserInfoRaw(
		"yingming",
		"男",
		"007",
		[]byte("123456789"),
		"http://mises.com",
		[]string{"yingming@gmail.com", "51911267@qq.com"},
		[]string{"17701314608", "18601350799", "18811790787"},
		"",
	)

	sessionid, err := user.SetUInfo(cuser, *info)
	assert.NoError(t, err)
	assert.False(t, sessionid == "")

	fmt.Printf("userinfo sessionid is %s\n", sessionid)
	PollSession(t, sessionid)

}

func TestFollow(t *testing.T) {
	cuser1 := CreateUserTest()

	PrepareUser(t, cuser1)

	cuser2 := CreateUserTest()

	PrepareUser(t, cuser2)

	sessionid, err := user.SetFollowing(cuser1, cuser2.MisesID(), "follow")
	assert.NoError(t, err)
	assert.False(t, sessionid == "")

	fmt.Printf("follow sessionid is %s\n", sessionid)
	PollSession(t, sessionid)

}

func TestParseTxResp(t *testing.T) {

	resp := user.MsgTxResp{}
	resp.Code = 0
	resp.Error = ""
	resp.TxResponse = user.MsgTx{Height: "1", Txhash: "123456"}
	respBytes, err := json.Marshal(resp)
	fmt.Printf("resp is %s\n", string(respBytes))
	assert.NoError(t, err)
	msgTx, err := user.ParseTxResp(respBytes)
	assert.NoError(t, err)
	assert.EqualString(t, msgTx.Txhash, "123456")
}
