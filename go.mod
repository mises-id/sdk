module github.com/mises-id/sdk

go 1.16

require (
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/btcutil v1.0.4
	github.com/cosmos/cosmos-sdk v0.44.5
	github.com/ebfe/keccak v0.0.0-20150115210727-5cc570678d1b
	github.com/gorilla/mux v1.8.0
	github.com/mises-id/mises-tm v0.0.0-20210821062909-5f9ffc470b61
	github.com/multiformats/go-multibase v0.0.3 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.6
	github.com/tyler-smith/assert v1.0.1
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3

)

replace google.golang.org/grpc => google.golang.org/grpc v1.42.0

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/tendermint/tm-db => github.com/mises-id/tm-db v0.6.5-0.20210822095222-e1ff1e0dc734

replace github.com/cosmos/iavl => github.com/mises-id/iavl v0.17.4-0.20211207035003-f9d26e6150db

replace github.com/tendermint/tendermint => github.com/mises-id/tendermint v0.34.15-0.20211207033151-1f29b59c0edf

replace github.com/mises-id/mises-tm => ../../core/mises-tm

replace github.com/cosmos/cosmos-sdk => github.com/mises-id/cosmos-sdk v0.44.6-0.20211209094558-a7c9c77cfc17
