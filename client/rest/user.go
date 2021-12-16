package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/mises-id/sdk/types"
)

type RestUserActiveRequest struct {
}

type RestUserActiveResponse struct {
	MisesId string
}

// HandleUserActiveRequest the UserActiveRequest http handler
func HandleUserActiveRequest(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		resp := &RestUserActiveResponse{}

		if KeyActivated != nil {
			resp.MisesId = types.MisesIDPrefix + KeyActivated.Address
		}

		PostProcessResponseBare(w, clientCtx, resp)
	}
}
