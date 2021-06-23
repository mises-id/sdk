package sdk

import (
	"os"
	"strconv"
	"time"
)

type test struct {
	SDK    *misesSdk
	Agent  *misesAgent
	Key1   string
	Key2   string
	Key3   string
	Value1 string
	Value2 string
	Value3 string
}

func (ctx *test) testSetUp() error {
	setupLogging()
	loadEnv()

	c, err := newTestSDK()
	if err != nil {
		return err
	} else {
		ctx.SDK = c.(*misesSdk)
	}

	ctx.Key1 = strconv.FormatInt(100+time.Now().Unix(), 10)
	ctx.Key2 = strconv.FormatInt(200+time.Now().Unix(), 10)
	ctx.Key3 = strconv.FormatInt(300+time.Now().Unix(), 10)

	ctx.Value1 = "foo"
	ctx.Value2 = "bar"
	ctx.Value3 = "baz"

	return nil
}

func (ctx *test) testTearDown() error {
	return nil
}

func newTestSDK() (MSdk, error) {
	debug := false
	if d, err := strconv.ParseBool(os.Getenv("DEBUG")); err == nil {
		debug = d
	}

	// create client
	options := &MSdkOption{
		ChainID: os.Getenv("CHAIN_ID"),
		Debug:   debug,
	}
	ctx, err := NewMSdk(options)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func setupLogging() {
}

func loadEnv() {
}

func testGasInfo() *GasInfo {
	return &GasInfo{
		MaxFee: 4000001,
	}
}

func testAddress() string {
	return os.Getenv("ADDRESS")
}
