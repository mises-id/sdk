package lcd

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	tmcfg "github.com/tendermint/tendermint/config"

	lproxy "github.com/tendermint/tendermint/light/proxy"

	"github.com/mises-id/sdk/client/cli/commands/light"
	"github.com/mises-id/sdk/types"
)

var _ MLightNode = &mLCD{}

type mLCD struct {
	chainId          string
	primaryAddress   string
	witnessAddresses string
	trustHeight      string
	trustHash        string
	insecureSsl      bool
	initThreadID     uint64
	proxy            *lproxy.Proxy
	restarting       uint32 // atomic
}

func (lcd *mLCD) SetChainID(chainId string) error {
	lcd.chainId = chainId
	return nil
}
func (lcd *mLCD) SetEndpoints(primary string, witnesses string) error {
	lcd.primaryAddress = primary
	lcd.witnessAddresses = witnesses
	return nil
}
func (lcd *mLCD) SetTrust(height string, hash string) error {
	lcd.trustHeight = height
	lcd.trustHash = hash
	return nil
}

func (lcd *mLCD) SetInsecureSsl(insecureSsl bool) error {
	lcd.insecureSsl = insecureSsl
	return nil
}

func (lcd *mLCD) serveImpl(listen string) error {
	_, err := CreateDefaultTendermintConfig(types.NodeHome)
	if err != nil {
		return err
	}

	trustHeight, err := strconv.ParseInt(lcd.trustHeight, 10, 64)
	if err != nil {
		trustHeight = 0
	}
	trustHash, err := hex.DecodeString(lcd.trustHash)
	if err != nil {
		trustHash = []byte{}
	}

	config := light.ProxyConfig{
		LogLevel:           "trace",
		ChainID:            lcd.chainId,
		WitnessAddrsJoined: lcd.witnessAddresses,
		Dir:                types.NodeHome + "/light",
		PrimaryAddr:        lcd.primaryAddress,
		TrustLevel:         "1/3",
		Sequential:         false,
		InsecureSsl:        lcd.insecureSsl,
		TrustedHeight:      trustHeight,
		TrustedHash:        trustHash,
		TrustingPeriod:     168 * time.Hour,
		MaxOpenConnections: 900,
		ListenAddr:         listen,
	}

	p, err := light.CreateProxy(&config)
	if err != nil {
		return err
	}
	lcd.proxy = p
	atomic.StoreUint32(&lcd.restarting, 0)

	if err := p.ListenAndServe(); err != http.ErrServerClosed {
		light.ClearProxy(p)
	}

	return nil
}

func (lcd *mLCD) Serve(listen string, delegator MLightNodeDelegator) error {

	go func() {
		for {
			lcd.serveImpl(listen)
			lcd.proxy = nil
			if atomic.LoadUint32(&lcd.restarting) == 1 {
				//restarting
				time.Sleep(5 * time.Second)
			} else {
				delegator.OnError()
				break
			}
		}
	}()

	return nil

}

func (lcd *mLCD) SetLogLevel(level int) error {
	return nil
}

func (lcd *mLCD) Restart() error {
	if atomic.CompareAndSwapUint32(&lcd.restarting, 0, 1) {
		proxy := lcd.proxy
		if proxy != nil {
			light.ClearProxy(proxy)
		}

	}
	return nil
}
func CreateDefaultTendermintConfig(rootDir string) (*tmcfg.Config, error) {
	conf := tmcfg.DefaultConfig()
	conf.SetRoot(rootDir)
	tmcfg.EnsureRoot(rootDir)

	if err := conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}

	return conf, nil
}

func NewMLightNode() MLightNode {
	lcd := &mLCD{}
	runtime.LockOSThread()
	return lcd
}

func SetHomePath(dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	types.NodeHome = dir + ".misestm"
	return nil
}
