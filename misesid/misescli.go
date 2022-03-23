package misesid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mises-id/sdk/types"
)

type CallBack func(session string, body []byte) (*WaitResult, error)
type WaitTask struct {
	session   string
	pCallback CallBack
}
type WaitResult struct {
	Session string
	Result  string
	ErrMsg  string
}

var (
	wrChan  = make(chan *WaitResult)
	APIHost = types.DefaultEndpoint
)

const (
	MisesIDURLPath = "mises/did"

	UInfoURLPath = "mises/user"

	FollowingURLPath = "mises/user/relation"

	TxURLPath = "mises/tx"
)

func MakeGetUrl(urlPath string, queryParams string, signer types.MSigner) (string, error) {
	signerMisesID := signer.MisesID()

	nonce := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	queryRequest := urlPath + "?" + queryParams + "&nonce=" + nonce + "&signer=" + signerMisesID
	signed, err := signer.Sign(queryRequest)
	if err != nil {
		return "", err
	}

	url := APIHost + queryRequest + "&sig=" + signed
	return url, nil
}

func HttpGetTx(sessionID string, cuser types.MSigner) ([]byte, error) {
	url, err := MakeGetUrl(TxURLPath, "txhash="+sessionID, cuser)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
func WaitResp(t WaitTask, cuser types.MSigner) {

	for i := 0; i < 12; i++ {
		time.Sleep(time.Second * 5)

		body, err := HttpGetTx(t.session, cuser)
		if err != nil {
			continue
		}

		ret, err := t.pCallback(t.session, body)
		if err != nil {
			continue
		}
		wrChan <- ret

		return
	}

	var r WaitResult
	r.Session = t.session
	r.ErrMsg = "timout"
	wrChan <- &r
}

func BuildPostForm(msg interface{}, cuser types.MSigner) (url.Values, error) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	msgStr := string(msgBytes)
	nonce := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	sigData := url.Values{}
	sigData.Set("msg", msgStr)
	sigData.Set("nonce", nonce)
	signed, err := cuser.Sign(sigData.Encode())
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("msg", msgStr)
	v.Set("nonce", nonce)
	v.Set("sig", signed)
	return v, nil
}

// Retry update to frontend
func CreateMisesID(signer types.MSigner) (string, error) {
	msg := MsgCreateMisesID{
		MsgReqBase: MsgReqBase{signer.MisesID()},
		PubKey:     signer.PubKey(),
	}
	v, err := BuildPostForm(&msg, signer)
	if err != nil {
		return "", err
	}

	return Set2Mises(signer, APIHost+MisesIDURLPath, v)
}
func GetMisesID(cuser types.MSigner, misesid string) (string, error) {
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
	return ParseGetMisesIDResp(body)
}

// retry up to frontend
func GetUInfo(cuser types.MSigner, misesid string) (*MisesUserInfo, error) {
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
	respMsg, err := ParseGetUserInfoResp(body)
	if err != nil {
		return nil, err
	}
	uinfoByte, err := Decrypt(cuser, respMsg.PrivateInfo.EncData, respMsg.PrivateInfo.IV)
	if err != nil {
		return nil, err
	}
	uinfo := MisesUserInfo{}
	err = json.Unmarshal(uinfoByte, &uinfo)
	if err != nil {
		return nil, err
	}

	return &uinfo, nil

}

func GetFollowing(cuser types.MSigner, misesid string) ([]string, error) {
	url, err := MakeGetUrl(FollowingURLPath, "filter=following&mises_id="+misesid, cuser)
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
	msgFolllowList, err := ParseListMisesResp(body)
	if err != nil {
		return nil, err
	}
	mids := []string{}
	for _, following := range msgFolllowList {
		mids = append(mids, following.MisesId)
	}

	return mids, nil
}

func SetUInfo(cuser types.MSigner, uinfo *MisesUserInfo) (string, error) {
	uinfoByte, err := json.Marshal(uinfo)
	if err != nil {
		return "", err
	}
	enc, iv, err := Encrypt(cuser, uinfoByte)
	if err != nil {
		return "", err
	}
	encData := EncryptedData{
		EncData: enc,
		IV:      iv,
	}
	msg := MsgUpdateUserInfo{
		MsgReqBase:  MsgReqBase{cuser.MisesID()},
		PublicInfo:  *uinfo,
		PrivateInfo: encData,
	}
	v, err := BuildPostForm(&msg, cuser)
	if err != nil {
		return "", err
	}

	return Set2Mises(cuser, APIHost+UInfoURLPath, v)
}

func SetFollowing(cuser types.MSigner, followingId string, op string) (string, error) {
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

func Set2Mises(cuser types.MSigner, url string, v url.Values) (string, error) {
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
	if r.Code != 0 || r.Error != "" {
		return nil, fmt.Errorf("failed to query tx:" + r.Error)
	}

	return &r.TxResponse, nil
}

func ParseListMisesResp(body []byte) ([]MsgMises, error) {
	var r MsgListMisesResp

	fmt.Println("ParseListMisesResp " + string(body))
	err := json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("failed to list mises:" + r.Error)
	}

	return r.MisesList, nil
}

func ParseGetUserInfoResp(body []byte) (*MsgGetUserInfoResp, error) {
	var r MsgGetUserInfoResp

	fmt.Println("ParseGetUserInfoResp " + string(body))
	err := json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("failed to get uinfo:" + r.Error)
	}

	return &r, nil
}
func ParseGetMisesIDResp(body []byte) (string, error) {
	var r MsgGetMisesIDResp

	err := json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}
	if r.Code != 0 || len(r.Error) > 0 {
		return "", fmt.Errorf("failed to get mises id:" + r.Error)
	}

	return r.PubKey, nil
}

func QueryCallBack(session string, body []byte) (*WaitResult, error) {
	var r = WaitResult{
		Session: session,
	}

	MsgTx, err := ParseTxResp(body)
	if err != nil {
		return nil, err
	}
	if MsgTx.Txhash == r.Session {
		r.Result = MsgTx.Height
		if MsgTx.Code != 0 {
			r.ErrMsg = MsgTx.RawLog
		}
	}

	return &r, nil
}

func PollSessionResult(timeout time.Duration) (*WaitResult, error) {
	select {
	case ret := <-wrChan:
		return ret, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("no result")
	}

}
func SetTestEndpoint(endpoint string) error {
	APIHost = endpoint
	return nil
}
