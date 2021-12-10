package commands

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/tendermint/spm/openapiconsole"
	tmrpcserver "github.com/tendermint/tendermint/rpc/jsonrpc/server"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/log"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/mises-id/mises-tm/docs"
	"github.com/mises-id/mises-tm/x/misestm/client/rest"
	"github.com/mises-id/sdk/types"
)

// simd light cosmoshub-3 --primary-addr http://193.26.156.221:26657/ --witness-addr http://144.76.61.201:26657/ --trusted-height 5940895 --trusted-hash 8663FBD3FB9DCE3D8E461EA521C38256F6EAF85D4FA492BAE26D5863F53CA150

func RestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rest",
		Short:   "Run a rest server",
		Long:    `Run a rest server`,
		RunE:    runRest,
		Example: `rest`,
	}

	cmd.Flags().String(listenAddrOpt, "tcp://localhost:1317", "serve the proxy on the given address")
	cmd.Flags().String(logLevelOpt, "info", "Log level, info or debug (Default: info) ")

	cmd.Flags().Int(maxOpenConnectionsOpt, 900, "maximum number of simultaneous connections (including WebSocket).")

	cmd.PersistentFlags().String(flags.FlagChainID, "test", "The network chain ID")
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, "test", "keyring")
	cmd.PersistentFlags().String(flags.FlagNode, "tcp://localhost:26657", "local light node")

	cmd.PersistentFlags().String(flags.FlagHome, types.DefaultNodeHome, "home dir")

	return cmd
}

func runRest(cmd *cobra.Command, args []string) error {
	// Initialize logger.
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	var option log.Option
	logLevel, _ := cmd.Flags().GetString(logLevelOpt)
	if logLevel == "info" {
		option, _ = log.AllowLevel("info")
	} else {
		option, _ = log.AllowLevel("debug")
	}
	logger = log.NewFilter(logger, option)

	rtr := mux.NewRouter()
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return err
	}
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(codec, tx.DefaultSignModes)
	clientCtx = clientCtx.
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithCodec(codec).
		WithInterfaceRegistry(interfaceRegistry).
		WithTxConfig(txCfg)
	rest.RegisterRoutes(clientCtx, rtr, true)
	rtr.Handle("/static/mises.yml", http.FileServer(http.FS(docs.Docs)))
	rtr.HandleFunc("/", openapiconsole.Handler("mises light", "/static/mises.yml"))

	tmCfg := tmrpcserver.DefaultConfig()
	maxOpenConnections, err := cmd.Flags().GetInt(maxOpenConnectionsOpt)
	if err != nil {
		return err
	}

	tmCfg.MaxOpenConnections = maxOpenConnections

	listenAddr, _ := cmd.Flags().GetString(listenAddrOpt)

	listener, err := tmrpcserver.Listen(listenAddr, tmCfg)
	if err != nil {
		return err
	}
	return tmrpcserver.Serve(listener, rtr, logger, tmCfg)
}
