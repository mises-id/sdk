package rest

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/cosmos/cosmos-sdk/types/rest"
)

var (
	KeyringPass  = &PassReader{}
	KeyActivated *MisesKey
)

type PassReader struct {
	Pass string
}

func (r *PassReader) Read(p []byte) (n int, err error) {
	n = copy(p, []byte(r.Pass))
	n += copy(p[n:], []byte("\n"))
	return
}

type MisesKey struct {
	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	PubKey  string `protobuf:"bytes,2,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
	Address string `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
}

type RestKeysRequest struct {
	Name       string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Passphrase string `protobuf:"bytes,2,opt,name=passphrase,proto3" json:"passphrase,omitempty"`
	PriKey     string `protobuf:"bytes,3,opt,name=pri_key,json=pri_key,proto3" json:"pri_key,omitempty"`
	Sig        string `json:"sig,omitempty"`
}

type RestKeysResponse struct {
	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
}

type RestKeysListResponse struct {
	Keys       []*MisesKey         `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func PostProcessResponseBare(w http.ResponseWriter, ctx client.Context, body interface{}) {
	var (
		resp []byte
		err  error
	)

	switch b := body.(type) {
	case []byte:
		resp = b

	default:
		resp, err = json.Marshal(body)
		if rest.CheckInternalServerError(w, err) {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

// HandleKeysListRequest the KeysListRequest http handler
func HandleKeysListRequest(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		infos, err := clientCtx.Keyring.List()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		resp := &RestKeysListResponse{
			Keys: []*MisesKey{},
		}

		for _, info := range infos {
			resp.Keys = append(resp.Keys, &MisesKey{
				Name:    info.GetName(),
				PubKey:  info.GetPubKey().String(),
				Address: info.GetAddress().String(),
			})

		}
		PostProcessResponseBare(w, clientCtx, resp)
	}
}

// HandleKeysImportRequest the KeysImportRequest http handler
func HandleKeysImportRequest(clientCtx client.Context) http.HandlerFunc {
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
		decodedData, err := hex.DecodeString(req.PriKey)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if len(req.Passphrase) >= 8 {
			KeyringPass.Pass = req.Passphrase + "\n"
		}
		priv := &secp256k1.PrivKey{Key: decodedData}

		armored := crypto.EncryptArmorPrivKey(priv, KeyringPass.Pass, "")

		if err := clientCtx.Keyring.ImportPrivKey(req.Name, armored, KeyringPass.Pass); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		resp := &RestKeysResponse{
			Code: 0,
		}
		PostProcessResponseBare(w, clientCtx, resp)
	}
}

// HandleKeysDeleteRequest the KeysDeleteRequest http handler
func HandleKeysDeleteRequest(clientCtx client.Context) http.HandlerFunc {
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

		if err := clientCtx.Keyring.Delete(req.Name); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		resp := &RestKeysResponse{
			Code: 0,
		}
		PostProcessResponseBare(w, clientCtx, resp)
	}
}

// HandleKeysResetRequest the KeysResetRequest http handler
func HandleKeysResetRequest(clientCtx client.Context) http.HandlerFunc {
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

		if err := os.Remove(clientCtx.KeyringDir); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		KeyActivated = nil
		resp := &RestKeysResponse{
			Code: 0,
		}
		PostProcessResponseBare(w, clientCtx, resp)
	}
}

// HandleKeysActivateequest the KeysActivateequest http handler
func HandleKeysActivateRequest(clientCtx client.Context) http.HandlerFunc {
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

		if len(req.Passphrase) >= 8 {
			KeyringPass.Pass = req.Passphrase + "\n"
		}
		if len(req.Name) > 0 {
			info, err := clientCtx.Keyring.Key(req.Name)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			KeyActivated = &MisesKey{
				Name:    info.GetName(),
				PubKey:  info.GetPubKey().String(),
				Address: info.GetAddress().String(),
			}
			clientCtx = clientCtx.WithFromName(KeyActivated.Name).WithFromAddress(types.AccAddress(KeyActivated.Address))
		}

		resp := &RestKeysResponse{
			Code: 0,
		}
		PostProcessResponseBare(w, clientCtx, resp)
	}
}
