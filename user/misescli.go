package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type CallBack func(body []byte) (*WaitResult, error)
type WaitTask struct {
	session   string
	pCallback CallBack
}
type WaitResult struct {
	session string
	result  string
}

var wr map[string]*WaitResult

const (
	APIHost = "http://localhost:1317/"

	MisesIDURLPath   = "mises/did/"
	CreateMisesIDUrl = APIHost + MisesIDURLPath
	//QuryMisesIDUrl   = APIHost + MisesIDURLPath

	UInfoURLPath = "mises/user/"
	SetUInfoURL  = APIHost + UInfoURLPath
	//GetUInfURL   = APIHost + UInfoURLPath

	FollowingURLPath = "mises/user/relation/"
	SetFollowingURL  = APIHost + FollowingURLPath
	//GetFollowingURL  = APIHost + FollowingURLPath

	TxURLPath = "mises/tx/"
	//QueryTxResultURL = APIHost + TxURLPath
)

func MakeGetUrl(urlPath string, queryParams string, cuser MUser) (string, error) {
	signerMisesID := cuser.MisesID()
	queryRequest := urlPath + "?" + queryParams + "&" + "signer=" + signerMisesID
	signed, nonce, err := Sign(cuser, queryRequest)
	if err != nil {
		return "", err
	}

	url := APIHost + queryRequest + "&sig=" + signed + "&nonce=" + nonce
	return url, nil
}

func WaitResp(t WaitTask, cuser MUser) {
	var r WaitResult
	wr = make(map[string]*WaitResult)
	r.session = t.session

	url, err := MakeGetUrl(TxURLPath, "tx_hash="+t.session, cuser)
	if err != nil {
		r.result = "601 url error"
		wr[r.session] = &r
		return
	}

	for i := 0; i < 3; i++ {
		time.Sleep(5000 * time.Millisecond)
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

	r.result = err.Error()
	wr[r.session] = &r
}

func BuildPostForm(msg interface{}, cuser MUser) (url.Values, error) {
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
	return v, nil
}

// Retry update to frontend
func CreateUser(cuser MUser) (string, error) {
	msg := MsgCreateMisesID{
		MsgReqBase: MsgReqBase{cuser.MisesID()},
		PubKey:     cuser.PubKEY(),
	}
	v, err := BuildPostForm(&msg, cuser)
	if err != nil {
		return "", err
	}

	return Set2Mises(cuser, CreateMisesIDUrl, v)
}

// retry up to frontend
func GetUInfo(cuser MUser, misesid string) (*MisesUserInfo, error) {
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

func GetFollowing(cuser MUser, misesid string) ([]string, error) {
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

func SetUInfo(cuser MUser, uinfo MisesUserInfo) (string, error) {
	msg := MsgUpdateUserInfo{
		MsgReqBase: MsgReqBase{cuser.MisesID()},
		PublicInfo: uinfo,
	}
	v, err := BuildPostForm(&msg, cuser)
	if err != nil {
		return "", err
	}

	return Set2Mises(cuser, SetUInfoURL, v)
}

func SetFollowing(cuser MUser, followingId string, op string) (string, error) {
	msg := MsgFollowMisesID{
		MsgReqBase: MsgReqBase{cuser.MisesID()},
		TargetID:   followingId,
		Action:     op,
	}
	v, err := BuildPostForm(&msg, cuser)
	if err != nil {
		return "", err
	}

	return Set2Mises(cuser, SetFollowingURL, v)
}

func Set2Mises(cuser MUser, url string, v url.Values) (string, error) {
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

func QueryCallBack(body []byte) (*WaitResult, error) {
	var r WaitResult

	var qr MsgTxResp

	err := json.Unmarshal(body, &qr)
	if err != nil {
		return nil, err
	}

	r.result = qr.TxResponse.Txhash
	return &r, nil
}
