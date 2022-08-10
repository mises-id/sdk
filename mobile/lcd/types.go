package lcd

type MLightNode interface {
	SetChainID(chainId string) error
	SetEndpoints(primary string, witnesses string) error
	SetTrust(height string, hash string) error
	SetInsecureSsl(insecureSsl bool) error
	Serve(listen string, delegator MLightNodeDelegator) error
	SetLogLevel(level int) error
	Restart() error
}

type MLightNodeDelegator interface {
	OnError()
}
