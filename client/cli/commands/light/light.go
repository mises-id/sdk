package light

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/libs/log"
	tmmath "github.com/tendermint/tendermint/libs/math"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/light"
	lproxy "github.com/tendermint/tendermint/light/proxy"
	lrpc "github.com/tendermint/tendermint/light/rpc"
	dbs "github.com/tendermint/tendermint/light/store/db"
	rpcserver "github.com/tendermint/tendermint/rpc/jsonrpc/server"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	dbm "github.com/tendermint/tm-db"
	_ "github.com/tendermint/tm-db/metadb"
)

// simd light cosmoshub-3 --primary-addr http://193.26.156.221:26657/ --witness-addr http://144.76.61.201:26657/ --trusted-height 5940895 --trusted-hash 8663FBD3FB9DCE3D8E461EA521C38256F6EAF85D4FA492BAE26D5863F53CA150

const (
	listenAddrOpt         = "listening-address"
	primaryAddrOpt        = "primary-addr"
	witnessAddrsJoinedOpt = "witness-addr"
	dirOpt                = "dir"
	maxOpenConnectionsOpt = "max-open-connections"
	insecureSslOpt        = "insecure-ssl"

	sequentialOpt     = "sequential-verification"
	trustingPeriodOpt = "trust-period"
	trustedHeightOpt  = "trusted-height"
	trustedHashOpt    = "trusted-hash"
	trustLevelOpt     = "trust-level"

	logLevelOpt = "log-level"
)

var (
	primaryKey   = []byte("primary")
	witnessesKey = []byte("witnesses")
	logger       = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "light")
)

func LightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "light [chainID]",
		Short: "Run a light client proxy server, verifying Tendermint rpc",
		Long: `Run a light client proxy server, verifying Tendermint rpc.
All calls that can be tracked back to a block header by a proof
will be verified before passing them back to the caller. Other than
that, it will present the same interface as a full Tendermint node.
Furthermore to the chainID, a fresh instance of a light client will
need a primary RPC address, a trusted hash and height and witness RPC addresses
(if not using sequential verification). To restart the node, thereafter
only the chainID is required.
`,
		RunE: runProxy,
		Args: cobra.ExactArgs(1),
		Example: `light cosmoshub-4 -primary-addr http://52.57.29.196:26657 -witness-addr http://public-seed-node.cosmoshub.certus.one:26657
	--height 962118 --hash 28B97BE9F6DE51AC69F70E0B7BFD7E5C9CD1A595B7DC31AFF27C50D4948020CD`,
	}

	cmd.Flags().String(listenAddrOpt, "tcp://localhost:8888", "serve the proxy on the given address")
	cmd.Flags().String(primaryAddrOpt, "", "connect to a Tendermint node at this address")
	cmd.Flags().String(witnessAddrsJoinedOpt, "", "tendermint nodes to cross-check the primary node, comma-separated")
	cmd.Flags().String(dirOpt, os.ExpandEnv(filepath.Join("$HOME", ".tendermint-light")), "specify the directory")
	cmd.Flags().Int(maxOpenConnectionsOpt, 900, "maximum number of simultaneous connections (including WebSocket).")
	cmd.Flags().Duration(trustingPeriodOpt, 168*time.Hour, "trusting period that headers can be verified within. Should be significantly less than the unbonding period")
	cmd.Flags().Int64(trustedHeightOpt, 1, "Trusted header's height")
	cmd.Flags().BytesHex(trustedHashOpt, []byte{}, "Trusted header's hash")
	cmd.Flags().String(logLevelOpt, "info", "Log level, info or debug (Default: info) ")
	cmd.Flags().String(trustLevelOpt, "1/3", "trust level. Must be between 1/3 and 3/3")
	cmd.Flags().Bool(sequentialOpt, false, "sequential verification. Verify all headers sequentially as opposed to using skipping verification")
	cmd.Flags().Bool(insecureSslOpt, false, "insecure skip ssl verification for android device below 7.1.1")

	return cmd
}

type ProxyConfig struct {
	PrimaryAddr        string
	WitnessAddrsJoined string
	LogLevel           string
	ChainID            string
	Dir                string
	TrustLevel         string
	TrustedHeight      int64
	TrustedHash        []byte
	TrustingPeriod     time.Duration
	MaxOpenConnections int
	ListenAddr         string

	Sequential  bool
	InsecureSsl bool
}

type ProxyState struct {
	Proxy *lproxy.Proxy
	DB    *dbm.PrefixDB
}

func ClearProxy(ps *ProxyState) {
	logger.Info("Clear Proxy")
	if ps.Proxy != nil {
		p := ps.Proxy
		if p.Listener != nil {
			err := p.Listener.Close()
			if err != nil {
				logger.Error("proxy close listener fail", err)
			}
		}
		if p.Client != nil && p.Client.IsRunning() {
			err := p.Client.Stop()
			if err != nil {
				logger.Error("proxy stop client fail", err)
			}

		}
	}

	if ps.DB != nil {
		defer ps.DB.Close()
	}

}

func CreateProxy(config *ProxyConfig) (*ProxyState, error) {

	var option log.Option
	if config.LogLevel == "info" {
		option, _ = log.AllowLevel("info")
	} else {
		option, _ = log.AllowLevel("debug")
	}
	logger = log.NewFilter(logger, option)

	logger.Info("Creating client...", "chainID", config.ChainID)

	var witnessesAddrs []string
	witnessesAddrs = []string{}
	if config.WitnessAddrsJoined != "" {
		witnessesAddrs = strings.Split(config.WitnessAddrsJoined, ",")
	}

	trustLevel, err := tmmath.ParseFraction(config.TrustLevel)
	if err != nil {
		return nil, fmt.Errorf("can't parse trust level: %w", err)
	}

	lightDB, err := dbm.NewGoLevelDB("light-client-db", config.Dir)
	if err != nil {
		return nil, fmt.Errorf("can't create a db: %w", err)
	}

	// create a prefixed db on the chainID
	db := dbm.NewPrefixDB(lightDB, []byte(config.ChainID))

	if config.PrimaryAddr == "" { // check to see if we can start from an existing state
		var err error
		var primaryAddress string
		primaryAddress, witnessesAddrs, err = checkForExistingProviders(db)
		if err != nil {
			defer db.Close()
			return nil, fmt.Errorf("failed to retrieve primary or witness from db: %w", err)
		}
		if primaryAddress == "" {
			defer db.Close()
			return nil, errors.New("no primary address was provided nor found. Please provide a primary (using -p)." +
				" Run the command: tendermint light --help for more information")
		}
		config.PrimaryAddr = primaryAddress
	} else {
		err := saveProviders(db, config.PrimaryAddr, config.WitnessAddrsJoined)
		if err != nil {
			logger.Error("Unable to save primary and or witness addresses", "err", err)
		}
	}

	options := []light.Option{
		light.Logger(logger),
		light.ConfirmationFunction(func(action string) bool {
			fmt.Println(action)
			return true
		}),
	}

	if config.InsecureSsl {
		rpctypes.ForceTrustIsrgRootX1()
	}

	if config.Sequential {
		options = append(options, light.SequentialVerification())
	} else {
		options = append(options, light.SkippingVerification(trustLevel))
	}

	var c *light.Client

	if config.TrustedHeight > 0 && len(config.TrustedHash) > 0 { // fresh installation
		c, err = light.NewHTTPClient(
			context.Background(),
			config.ChainID,
			light.TrustOptions{
				Period: config.TrustingPeriod,
				Height: config.TrustedHeight,
				Hash:   config.TrustedHash,
			},
			config.PrimaryAddr,
			witnessesAddrs,
			dbs.New(lightDB, config.ChainID),
			options...,
		)
	} else { // continue from latest state
		c, err = light.NewHTTPClientFromTrustedStore(
			config.ChainID,
			config.TrustingPeriod,
			config.PrimaryAddr,
			witnessesAddrs,
			dbs.New(lightDB, config.ChainID),
			options...,
		)
	}
	if err != nil {
		defer db.Close()
		return nil, err
	}

	tmconfig := tmcfg.DefaultConfig()

	cfg := rpcserver.DefaultConfig()
	cfg.MaxBodyBytes = tmconfig.RPC.MaxBodyBytes
	cfg.MaxHeaderBytes = tmconfig.RPC.MaxHeaderBytes

	cfg.MaxOpenConnections = config.MaxOpenConnections
	// If necessary adjust global WriteTimeout to ensure it's greater than
	// TimeoutBroadcastTxCommit.
	// See https://github.com/tendermint/tendermint/issues/3435
	if cfg.WriteTimeout <= tmconfig.RPC.TimeoutBroadcastTxCommit {
		cfg.WriteTimeout = tmconfig.RPC.TimeoutBroadcastTxCommit + 1*time.Second
	}

	lp, err := lproxy.NewProxy(c, config.ListenAddr, config.PrimaryAddr, cfg, logger, lrpc.KeyPathFn(MerkleKeyPathFn()))

	if err != nil {
		defer db.Close()
		return nil, err
	}

	return &ProxyState{
		Proxy: lp,
		DB:    db,
	}, nil
}

func runProxy(cmd *cobra.Command, args []string) error {
	// Initialize logger.
	logLevel, _ := cmd.Flags().GetString(logLevelOpt)

	chainID := args[0]

	witnessAddrsJoined, _ := cmd.Flags().GetString(witnessAddrsJoinedOpt)

	dir, _ := cmd.Flags().GetString(dirOpt)

	primaryAddr, _ := cmd.Flags().GetString(primaryAddrOpt)

	tl, err := cmd.Flags().GetString(trustLevelOpt)
	if err != nil {
		return err
	}

	sequential, err := cmd.Flags().GetBool(sequentialOpt)
	if err != nil {
		return err
	}

	insecureSsl, err := cmd.Flags().GetBool(insecureSslOpt)
	if err != nil {
		return err
	}

	trustedHeight, err := cmd.Flags().GetInt64(trustedHeightOpt)
	if err != nil {
		return err
	}

	trustedHash, err := cmd.Flags().GetBytesHex(trustedHashOpt)
	if err != nil {
		return err
	}
	trustingPeriod, err := cmd.Flags().GetDuration(trustingPeriodOpt)
	if err != nil {
		return err
	}

	maxOpenConnections, err := cmd.Flags().GetInt(maxOpenConnectionsOpt)
	if err != nil {
		return err
	}

	listenAddr, _ := cmd.Flags().GetString(listenAddrOpt)

	config := ProxyConfig{
		LogLevel:           logLevel,
		ChainID:            chainID,
		WitnessAddrsJoined: witnessAddrsJoined,
		Dir:                dir,
		PrimaryAddr:        primaryAddr,
		TrustLevel:         tl,
		Sequential:         sequential,
		InsecureSsl:        insecureSsl,
		TrustedHeight:      trustedHeight,
		TrustedHash:        trustedHash,
		TrustingPeriod:     trustingPeriod,
		MaxOpenConnections: maxOpenConnections,
		ListenAddr:         listenAddr,
	}

	ps, err := CreateProxy(&config)
	if err != nil {
		return err
	}

	// Stop upon receiving SIGTERM or CTRL-C.
	tmos.TrapSignal(logger, func() {
	})

	logger.Info("Starting proxy...", "laddr", listenAddr)
	if err := ps.Proxy.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logger.Error("proxy ListenAndServe", "err", err)
		ClearProxy(ps)

	}

	return nil
}

func checkForExistingProviders(db dbm.DB) (string, []string, error) {
	primaryBytes, err := db.Get(primaryKey)
	if err != nil {
		return "", []string{""}, err
	}
	witnessesBytes, err := db.Get(witnessesKey)
	if err != nil {
		return "", []string{""}, err
	}
	witnessesAddrs := strings.Split(string(witnessesBytes), ",")
	return string(primaryBytes), witnessesAddrs, nil
}

func saveProviders(db dbm.DB, primaryAddr, witnessesAddrs string) error {
	err := db.Set(primaryKey, []byte(primaryAddr))
	if err != nil {
		return fmt.Errorf("failed to save primary provider: %w", err)
	}
	err = db.Set(witnessesKey, []byte(witnessesAddrs))
	if err != nil {
		return fmt.Errorf("failed to save witness providers: %w", err)
	}
	return nil
}

// DefaultMerkleKeyPathFn creates a function used to generate merkle key paths
// from a path string and a key. This is the default used by the cosmos SDK.
// This merkle key paths are required when verifying /abci_query calls
func MerkleKeyPathFn() lrpc.KeyPathFunc {
	// regexp for extracting store name from /abci_query path
	storeNameRegexp := regexp.MustCompile(`\/store\/(.+)\/key`)

	return func(path string, key []byte) (merkle.KeyPath, error) {
		matches := storeNameRegexp.FindStringSubmatch(path)
		if len(matches) != 2 {
			return nil, fmt.Errorf("can't find store name in %s using %s", path, storeNameRegexp)
		}
		storeName := matches[1]

		kp := merkle.KeyPath{}
		kp = kp.AppendKey([]byte(storeName), merkle.KeyEncodingURL)
		kp = kp.AppendKey(key, merkle.KeyEncodingURL)
		return kp, nil
	}
}
