package user

import (
	//"encoding/json"
	"fmt"
	"time"

	//"io/ioutil"
	//"net/http"

	//"strings"
	"testing"

	"github.com/mises-id/sdk/bip39"
	"github.com/tyler-smith/assert"
)

/*
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
*/
func CreateUserTest() MUser {
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

func TestGetUInfo(t *testing.T) {
	cuser := CreateUserTest()

	url, err := MakeGetUrl(cuser.MisesID(), GetUInfoUrl, cuser)
	assert.Nil(t, err)

	resp, err := http.Get(url)
	assert.Nil(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	fmt.Printf("body is %s\n", string(body))

	var uinfo MisesUserInfo
	err = json.Unmarshal(body, &uinfo)
	assert.Nil(t, err)

	ub, err := json.Marshal(&uinfo)
	assert.Nil(t, err)

	fmt.Printf("user info is %s\n", string(ub))
}


func TestGetFollowing(t *testing.T) {
	cuser := CreateUserTest()

	uf := cuser.MisesID()
	url, err := MakeGetUrl(uf, GetFollowingUrl, cuser)
	assert.Nil(t, err)

	resp, err := http.Get(url)
	assert.Nil(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

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
	assert.Nil(t, err)

	fmt.Printf("following sessionid is %s\n", sessionid)

	time.Sleep(10000 * time.Millisecond)
}
*/

func TestSetUserInfo(t *testing.T) {
	cuser := CreateUserTest()

	var info MisesUserInfo

	info.Name = "yingming"
	info.Gender = "男"
	info.AvatarId = "007"
	info.AvatarThumb = []byte("123456789")
	info.HomePage = "http://mises.com"
	emails := []string{"yingming@gmail.com", "51911267@qq.com"}
	teles := []string{"17701314608", "18601350799", "18811790787"}
	info.Emails = emails
	info.Telephones = teles

	sessionid, err := SetUInfo(cuser, info)
	assert.Nil(t, err)

	fmt.Printf("userinfo sessionid is %s\n", sessionid)

	time.Sleep(10000 * time.Millisecond)
}
