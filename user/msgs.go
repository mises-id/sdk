package user

type MsgReqBase struct {
	MisesID string `json:"mises_id,omitempty"`
}

type MsgCreateMisesID struct {
	MsgReqBase
	PubKey string `json:"pub_key,omitempty"`
}

type EncryptedData struct {
	EncData string `json:"enc_data,omitempty"`
	IV      string `json:"iv,omitempty"`
}
type MsgUpdateUserInfo struct {
	MsgReqBase
	PublicInfo  MisesUserInfo `json:"pub_info,omitempty"`
	PrivateInfo EncryptedData `json:"pri_info,omitempty"`
}

type MsgFollowMisesID struct {
	MsgReqBase
	TargetID string `json:"target_id,omitempty"`
	Action   string `json:"action,omitempty"`
}

type MsgRespBase struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}
type MsgTx struct {
	Height uint64 `json:"height,omitempty"`
	Txhash string `json:"txhash,omitempty"`
}

type MsgTxResp struct {
	MsgRespBase
	TxResponse MsgTx `json:"tx_response,omitempty"`
}

type MsgPagination struct {
	NextKey string `json:"nextKey,omitempty"`
	Total   uint64 `json:"total,omitempty"`
}

type MsgFollow struct {
	MisesId string `json:"mises_id,omitempty"`
}

type MsgListFollowResp struct {
	MsgRespBase
	FollowList []MsgFollow   `json:"follow_list,omitempty"`
	Pagination MsgPagination `json:"pagination,omitempty"`
}

type MsgGetUserInfoResp struct {
	MsgRespBase
	PublicInfo  MisesUserInfo `json:"pub_info,omitempty"`
	PrivateInfo EncryptedData `json:"pri_info,omitempty"`
}
