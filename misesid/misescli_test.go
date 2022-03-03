package misesid_test

import (
	//"encoding/json"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	//"io/ioutil"
	//"net/http"

	//"strings"
	"testing"

	"github.com/mises-id/sdk/bip39"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sdk/user"
	"github.com/tyler-smith/assert"
)

func init() {
	/* load test data */
	misesid.SetTestEndpoint("http://gw.mises.site:1317/")
}

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
	wr, err := misesid.PollSessionResult(60 * time.Second)
	fmt.Printf("PollSessionResult finish\n")
	assert.NoError(t, err)
	assert.EqualString(t, "", wr.ErrMsg)
	assert.EqualString(t, session, wr.Session)
}
func PrepareUser(t *testing.T, cuser types.MUser) {
	sessionid, err := misesid.CreateMisesID(cuser.Signer())
	assert.NoError(t, err)
	assert.False(t, sessionid == "")
	fmt.Printf("create misesid sessionid is %s\n", sessionid)
	PollSession(t, sessionid)
}

func TestUserCreateMisesID(t *testing.T) {
	cuser := CreateUserTest()

	PrepareUser(t, cuser)

}

func TestUserSetInfo(t *testing.T) {
	cuser := CreateUserTest()

	PrepareUser(t, cuser)

	info := misesid.MisesUserInfo{
		"yingming",
		"male",
		"ipfs://asdasdasdadsa",
		"http://mises.com",
		[]string{"yingming@gmail.com", "51911267@qq.com"},
		[]string{"17701314608", "18601350799", "18811790787"},
		"",
	}

	sessionid, err := misesid.SetUInfo(cuser.Signer(), &info)
	assert.NoError(t, err)
	assert.False(t, sessionid == "")

	fmt.Printf("userinfo sessionid is %s\n", sessionid)
	PollSession(t, sessionid)

	respInfo, err := misesid.GetUInfo(cuser.Signer(), cuser.MisesID())
	assert.NoError(t, err)
	assert.EqualString(t, "yingming", respInfo.Name)

}

func TestUserFollow(t *testing.T) {
	cuser1 := CreateUserTest()

	PrepareUser(t, cuser1)

	cuser2 := CreateUserTest()

	PrepareUser(t, cuser2)

	sessionid, err := misesid.SetFollowing(cuser1.Signer(), cuser2.MisesID(), "follow")
	assert.NoError(t, err)
	assert.False(t, sessionid == "")

	fmt.Printf("follow sessionid is %s\n", sessionid)
	PollSession(t, sessionid)

	followingIDs, err := misesid.GetFollowing(cuser1.Signer(), cuser1.MisesID())
	assert.NoError(t, err)
	assert.True(t, followingIDs[0] == cuser2.MisesID())

	sessionid1, err := misesid.SetFollowing(cuser1.Signer(), cuser2.MisesID(), "unfollow")
	assert.NoError(t, err)
	assert.False(t, sessionid1 == "")

	fmt.Printf("unfollow sessionid is %s\n", sessionid1)
	PollSession(t, sessionid1)

	followingIDs1, err := misesid.GetFollowing(cuser1.Signer(), cuser1.MisesID())
	assert.NoError(t, err)
	assert.True(t, len(followingIDs1) == 0)

}

func TestUserParseTxResp(t *testing.T) {

	resp := misesid.MsgTxResp{}
	resp.Code = 0
	resp.Error = ""
	resp.TxResponse = misesid.MsgTx{Height: "1", Txhash: "123456"}
	respBytes, err := json.Marshal(&resp)
	fmt.Printf("resp is %s\n", string(respBytes))
	assert.NoError(t, err)
	msgTx, err := misesid.ParseTxResp(respBytes)
	assert.NoError(t, err)
	assert.EqualString(t, "123456", msgTx.Txhash)
}

func dummyUpdate(t *testing.T, cuser types.MUser) {
	//sessionid, _ := user.SetUInfo(cuser, &user.MisesUserInfo{})
	encData := misesid.EncryptedData{
		EncData: "ipfRvOlodErWniY/E+hHUTSn7yiw2PzOvXceQk0RsutToZIxBW+w+yDSzEI9A/1qsmhh4PPcpVzzG6eKH8mkhfajBGi7CQvLTFNjqMVeJos=",
		IV:      "gONDIeRF2LNrq7vVDC/YXw==",
	}
	msg := misesid.MsgUpdateUserInfo{
		MsgReqBase:  misesid.MsgReqBase{cuser.MisesID()},
		PrivateInfo: encData,
	}
	v, _ := misesid.BuildPostForm(&msg, cuser.Signer())
	misesid.Set2Mises(cuser.Signer(), misesid.APIHost+misesid.UInfoURLPath, v)
	//PollSession(t, sessionid)
}
func TestGas(t *testing.T) {
	//user.SetTestEndpoint("http://gw.mises.site:1317/")
	tx1, _ := base64.StdEncoding.DecodeString("Cp4DCpUDCiovbWlzZXNpZC5taXNlc3RtLm1pc2VzdG0uTXNnVXBkYXRlVXNlckluZm8S5gIKLG1pc2VzMWc3YWhoOXY0dDVodDBkenp4dGc5bTlycWp5cWdlcHp0NGZjZ2pqEP///////////wEaNmRpZDptaXNlczptaXNlczE4czBrZnBtcHR5cHF4aDl4c2NyZW4wZ3lzN25ycHh1aDVldGF2eCLYAXBaRUhsaXhTMjk0Wm5kU1dkRExabGhCbmRLZXJVODBNYVRyNjZDVldjQ2JIWU9nM1o5UUlacUNuNTZMdHNrbzhOTjU3eXBOZlZlMmFpNkIrclFiMXptR0ZZdVNsUkg3WU1wSUVrS3FkdUQzMmg1VWtaQVU3bnhEWm1CdmNCbWV6MVpXNzEwWkgyN1JyNEhCOGMyQ0JWQm9xcHJuN1I5Rzd4T29iSGRCZUt3dit2T0dJaGxUY0cydzFYNk1LVEx4WnZKRVRoeDBtbnlzM2JRK05rUHZhT3c9PSoYSVJZbDZ2RFdXMzJjRFBRcGJXM2pZUT09EgRtZW1vElgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQNqXl6rNcwnDTt+EGC8kVC8LTj6TswsC8cni/wK4YqDTRIECgIIARgkEgQQoI0GGkBTvPnZerC1v7J6+2trbZQX9/bZqLB1zQQsmVxb/NxOE0rUAUjtiSkFWOAJPJxikBEr+AB6bx7DKu+qFgIxXGug")
	tx2, _ := base64.StdEncoding.DecodeString("Cp4DCpUDCiovbWlzZXNpZC5taXNlc3RtLm1pc2VzdG0uTXNnVXBkYXRlVXNlckluZm8S5gIKLG1pc2VzMWc3YWhoOXY0dDVodDBkenp4dGc5bTlycWp5cWdlcHp0NGZjZ2pqEP///////////wEaNmRpZDptaXNlczptaXNlczE3ajhwbGU2Z2N3eTVtYTY5a24zazVtdzB0NXpqMnY4NmV5bnltOCLYAXBaRUhsaXhTMjk0Wm5kU1dkRExabGhCbmRLZXJVODBNYVRyNjZDVldjQ2JIWU9nM1o5UUlacUNuNTZMdHNrbzhOTjU3eXBOZlZlMmFpNkIrclFiMXptR0ZZdVNsUkg3WU1wSUVrS3FkdUQzMmg1VWtaQVU3bnhEWm1CdmNCbWV6MVpXNzEwWkgyN1JyNEhCOGMyQ0JWQm9xcHJuN1I5Rzd4T29iSGRCZUt3dit2T0dJaGxUY0cydzFYNk1LVEx4WnZKRVRoeDBtbnlzM2JRK05rUHZhT3c9PSoYSVJZbDZ2RFdXMzJjRFBRcGJXM2pZUT09EgRtZW1vElgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQNqXl6rNcwnDTt+EGC8kVC8LTj6TswsC8cni/wK4YqDTRIECgIIARg4EgQQoI0GGkAGkHY713MGma0y+nQNBqujSOBDc4Ed+3dyJiBJzzkWdRBcCDuaUPYzy04anOzE2szR90UdX/gMeekeNZgg/FAN")
	fmt.Printf("len tx1 is %d\n", len(tx1))
	fmt.Printf("len tx2 is %d\n", len(tx2))
	cuser := CreateUserTest()

	PrepareUser(t, cuser)
	dummyUpdate(t, cuser)
	dummyUpdate(t, cuser)
	dummyUpdate(t, cuser)
	dummyUpdate(t, cuser)
	dummyUpdate(t, cuser)
	encData := misesid.EncryptedData{
		EncData: "pZEHlixS294ZndSWdDLZlhBndKerU80MaTr66CVWcCbHYOg3Z9QIZqCn56Ltsko8NN57ypNfVe2ai6B+rQb1zmGFYuSlRH7YMpIEkKqduD32h5UkZAU7nxDZmBvcBmez1ZW710ZH27Rr4HB8c2CBVBoqprn7R9G7xOobHdBeKwv+vOGIhlTcG2w1X6MKTLxZvJEThx0mnys3bQ+NkPvaOw==",
		IV:      "IRYl6vDWW32cDPQpbW3jYQ==",
	}
	msg := misesid.MsgUpdateUserInfo{
		MsgReqBase:  misesid.MsgReqBase{cuser.MisesID()},
		PrivateInfo: encData,
	}
	v, err := misesid.BuildPostForm(&msg, cuser.Signer())
	sessionid, err := misesid.Set2Mises(cuser.Signer(), misesid.APIHost+misesid.UInfoURLPath, v)

	assert.NoError(t, err)
	assert.False(t, sessionid == "")

	fmt.Printf("userinfo sessionid is %s\n", sessionid)
	PollSession(t, sessionid)

}
