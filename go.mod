module github.com/mises-id/sdk

go 1.16

require (
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/cosmos/btcutil v1.0.4
	github.com/cosmos/cosmos-sdk v0.44.6
	github.com/ebfe/keccak v0.0.0-20150115210727-5cc570678d1b
	github.com/mises-id/mises-tm v0.0.0-20220303064252-ef3c1ed6ee27
	github.com/multiformats/go-multibase v0.0.3
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/tendermint/tendermint v0.34.16
	github.com/tendermint/tm-db v0.6.6
	github.com/tyler-smith/assert v1.0.1
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3

)

replace google.golang.org/grpc => google.golang.org/grpc v1.42.0

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/tendermint/tm-db => github.com/mises-id/tm-db v0.6.5-0.20210822095222-e1ff1e0dc734

replace github.com/cosmos/iavl => github.com/mises-id/iavl v0.17.4-0.20211207035003-f9d26e6150db

replace github.com/cosmos/cosmos-sdk => github.com/mises-id/cosmos-sdk v0.44.6-0.20220315093538-763383563639

replace github.com/tendermint/tendermint => github.com/mises-id/tendermint v0.34.15-0.20220725013722-fd06dc4fa7a8

replace github.com/99designs/keyring => github.com/99designs/keyring v1.2.1
