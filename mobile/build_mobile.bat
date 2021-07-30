set GO111MODULE=
set GOMOD=
go env -w GO111MODULE=auto
gomobile bind -v  -o sdk.aar -target=android ./mobile