package user

import (
	"encoding/hex"
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

var wr map[string]WaitResult

var CreateUrl = "http://localhost:8080/create"
var QueryResultUrl = "http://localhost:8080/query?"
var GetUInfoUrl = "http://localhost:8080/uinfo?"
var GetFollowingUrl = "http://localhost:8080/following?"

//var GetFollowerUrl = "https://node.mises.site/follower?"
var SetUInfoUrl = "http://localhost:8080/setuinfo"
var SetFollowingUrl = "http://localhost:8080/setfollowing"

func MakeGetUrl(session string, baseUrl string, cuser MUser) (string, error) {
	msg, signed, err := Sign(cuser, session)
	if err != nil {
		return "", err
	}

	url := baseUrl + "msg=" + msg + "&sig=" + signed
	return url, nil
}

func WaitResp(t WaitTask, cuser MUser) {
	var r WaitResult
	r.session = t.session

	url, err := MakeGetUrl(t.session, QueryResultUrl, cuser)
	if err != nil {
		r.result = "601 url error"
		wr[r.session] = r
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
		wr[r.session] = *ret

		return
	}

	r.result = err.Error()
	wr[r.session] = r
}

// Retry update to frontend
func CreateUser(cuser MUser) (string, error) {
	msg, signed, err := Sign(cuser, cuser.MisesID())
	if err != nil {
		return "", err
	}

	fmt.Printf("misesid is: %s\n", cuser.MisesID())
	fmt.Printf("signed is %s\n", signed)

	v := url.Values{}
	v.Set("msg", msg)
	v.Set("sig", signed)

	return Set2Mises(cuser, CreateUrl, v)
}

// retry up to frontend
func GetUInfo(cuser MUser, misesid string) ([]byte, error) {
	url, err := MakeGetUrl(misesid, GetUInfoUrl, cuser)
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

	return body, nil
}

func GetFollowing(cuser MUser, misesid string) ([]byte, error) {
	url, err := MakeGetUrl(misesid, GetFollowingUrl, cuser)
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

	return body, nil
}

func SetUInfo(cuser MUser, uinfo MisesUserInfo) (string, error) {
	content, err := json.Marshal(uinfo)
	if err != nil {
		return "", err
	}

	info := hex.EncodeToString(content)
	m := cuser.MisesID() + " " + info
	msg, signed, err := Sign(cuser, m)
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Set("msg", msg)
	v.Set("sig", signed)

	return Set2Mises(cuser, SetUInfoUrl, v)
}

func SetFollowing(cuser MUser, followingId string, op string) (string, error) {
	m := cuser.MisesID() + " " + followingId + " " + op
	msg, signed, err := Sign(cuser, m)
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Set("msg", msg)
	v.Set("sig", signed)

	return Set2Mises(cuser, CreateUrl, v)
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
	t.session, err = ParseResp(body)
	if err != nil {
		return "", err
	}

	t.pCallback = QueryCallBack
	go WaitResp(t, cuser)

	return t.session, nil
}

func ParseResp(body []byte) (string, error) {
	/*
		var r string

		err := json.Unmarshal(body, &r)
		if err != nil {
			return "", err
		}

		return r, nil
	*/
	return string(body), nil
}

func QueryCallBack(body []byte) (*WaitResult, error) {
	var r WaitResult
	var qr string

	err := json.Unmarshal(body, &qr)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
