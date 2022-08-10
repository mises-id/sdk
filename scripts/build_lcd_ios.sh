###
 # @Author: lmk
 # @Date: 2021-07-30 22:55:01
 # @LastEditTime: 2021-07-31 17:25:27
 # @LastEditors: lmk
 # @Description: 
### 
export GO111MODULE=auto
export GOMOD=auto
export CGO_ENABLED=1
go env -w GO111MODULE=auto
gomobile bind -v  -o sdk.xcframework -target=ios ./mobile/lcd