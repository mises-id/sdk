package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/mises-id/sdk/types"
)

type CallBack func(body []byte) (*WaitResult, error)
type WaitTask struct {
	session   string
	pCallback CallBack
}
type WaitResult struct {
	session string
	result  string
	err     string
}

var (
	wr      = map[string]*WaitResult{}
	APIHost = types.DefaultEndpoint
)

const (
	MisesIDURLPath = "mises/did"

	UInfoURLPath = "mises/user"

	FollowingURLPath = "mises/user/relation"

	TxURLPath = "mises/tx"
)

func MakeGetUrl(urlPath string, queryParams string, cuser types.MUser) (string, error) {
	signerMisesID := cuser.MisesID()
	queryRequest := urlPath + "?" + queryParams + "&" + "signer=" + signerMisesID
	signed, nonce, err := Sign(cuser, queryRequest)
	if err != nil {
		return "", err
	}

	url := APIHost + queryRequest + "&sig=" + signed + "&nonce=" + nonce
	return url, nil
}

func WaitResp(t WaitTask, cuser types.MUser) {
	var r WaitResult
	r.session = t.session

	url, err := MakeGetUrl(TxURLPath, "tx_hash="+t.session, cuser)
	if err != nil {
		r.err = "601 url error"
		wr[r.session] = &r
		return
	}

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * time.Duration(2*i+1))
		resp, err := http.Get(url)
		if err != nil {
			continue
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		ret, err := t.pCallback(body)
		if err != nil {
			continue
		}
		wr[r.session] = ret

		return
	}

	r.err = err.Error()
	wr[r.session] = &r
}

func BuildPostForm(msg interface{}, cuser types.MUser) (url.Values, error) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	msgStr := string(msgBytes)
	signed, nonce, err := Sign(cuser, msgStr)
	if err != nil {
		return nil, err
	}
	fmt.Printf("msg is: %s\n", msgStr)
	fmt.Printf("signed is %s\n", signed)

	v := url.Values{}
	v.Set("msg", msgStr)
	v.Set("nonce", nonce)
	v.Set("sig", signed)
	fmt.Printf("post form is: %s\n", v.Encode())
	return v, nil
}

// Retry update to frontend
func CreateUser(cuser types.MUser) (string, error) {
	msg := MsgCreateMisesID{
		MsgReqBase: MsgReqBase{cuser.MisesID()},
		PubKey:     cuser.PubKEY(),
	}
	v, err := BuildPostForm(&msg, cuser)
	if err != nil {
		return "", err
	}

	return Set2Mises(cuser, APIHost+MisesIDURLPath, v)
}
func GetUser(cuser types.MUser, misesid string) (string, error) {
	url, err := MakeGetUrl(MisesIDURLPath, "mises_id="+misesid, cuser)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return ParseGetUserResp(body)
}

// retry up to frontend
func GetUInfo(cuser types.MUser, misesid string) (*MisesUserInfo, error) {
	url, err := MakeGetUrl(UInfoURLPath, "mises_id="+misesid, cuser)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ParseGetUserInfoResp(body)
}

func GetFollowing(cuser types.MUser, misesid string) ([]string, error) {
	url, err := MakeGetUrl(FollowingURLPath, "mises_id="+misesid, cuser)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	msgFolllowList, err := ParseListFollowResp(body)
	if err != nil {
		return nil, err
	}
	mids := []string{}
	for _, following := range msgFolllowList {
		mids = append(mids, following.MisesId)
	}

	return mids, nil
}

func SetUInfo(cuser types.MUser, uinfo MisesUserInfo) (string, error) {
	msg := MsgUpdateUserInfo{
		MsgReqBase: MsgReqBase{cuser.MisesID()},
		PublicInfo: uinfo,
	}
	v, err := BuildPostForm(&msg, cuser)
	if err != nil {
		return "", err
	}

	return Set2Mises(cuser, APIHost+UInfoURLPath, v)
}

func SetFollowing(cuser types.MUser, followingId string, op string) (string, error) {
	msg := MsgFollowMisesID{
		MsgReqBase: MsgReqBase{cuser.MisesID()},
		TargetID:   followingId,
		Action:     op,
	}
	v, err := BuildPostForm(&msg, cuser)
	if err != nil {
		return "", err
	}

	return Set2Mises(cuser, APIHost+FollowingURLPath, v)
}

func Set2Mises(cuser types.MUser, url string, v url.Values) (string, error) {
	fmt.Printf("post url is: %s\n", url)
	resp, err := http.PostForm(url, v)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var t WaitTask
	MsgTx, err := ParseTxResp(body)
	if err != nil {
		return "", err
	}
	t.session = MsgTx.Txhash

	t.pCallback = QueryCallBack
	go WaitResp(t, cuser)

	return t.session, nil
}

func ParseTxResp(body []byte) (*MsgTx, error) {

	var r MsgTxResp

	fmt.Println("ParseTxResp " + string(body))
	err := json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("failed to query tx:" + r.Error)
	}

	return &r.TxResponse, nil
}

func ParseListFollowResp(body []byte) ([]MsgFollow, error) {
	var r MsgListFollowResp

	err := json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("failed to list following:" + r.Error)
	}

	return r.FollowList, nil
}

func ParseGetUserInfoResp(body []byte) (*MisesUserInfo, error) {
	var r MsgGetUserInfoResp

	err := json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("failed to get uinfo:" + r.Error)
	}

	return &r.PublicInfo, nil
}
func ParseGetUserResp(body []byte) (string, error) {
	var r MsgGetUserResp

	err := json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}
	if r.Code != 0 {
		return "", fmt.Errorf("failed to get user:" + r.Error)
	}

	return r.PubKey, nil
}

func QueryCallBack(body []byte) (*WaitResult, error) {
	var r WaitResult

	var qr MsgTxResp

	err := json.Unmarshal(body, &qr)
	if err != nil {
		return nil, err
	}
	if qr.TxResponse.Txhash == r.session {
		r.result = qr.TxResponse.Height
	}

	return &r, nil
}

func CheckSession(sessinID string) (bool, error) {
	r, ok := wr[sessinID]
	if !ok {
		return false, fmt.Errorf("no such session " + sessinID)
	}
	if r.result != "0" {
		return true, nil
	}
	return false, nil

}
func SetTestEndpoint(endpoint string) error {
	APIHost = endpoint
	return nil
}
