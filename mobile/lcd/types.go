package lcd

type MLightNode interface {
	SetChainID(chainId string) error
	SetEndpoints(primary string, witnesses string) error
	SetTrust(height string, hash string) error
	Serve(listen string) error
	SetLogLevel(level int) error
}
