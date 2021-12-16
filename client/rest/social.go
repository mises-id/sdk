package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/types/rest"
)

type RestSocialForwardRequest struct {
	MisesID string `json:"mises_id"`
	Comment string `json:"comment"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	IconUrl string `json:"icon_url"`
}

type RestSocialForwardResponse struct {
	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
}

// HandleSocialForwardRequest the SocialForwardRequest http handler
func HandleSocialForwardRequest(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RestKeysRequest
		reqMsg, err := ioutil.ReadAll(r.Body)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		err = json.Unmarshal(reqMsg, &req)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		resp := &RestSocialForwardResponse{}

		PostProcessResponseBare(w, clientCtx, resp)
	}
}
