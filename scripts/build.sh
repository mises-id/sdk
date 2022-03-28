###
 # @Author: lmk
 # @Date: 2021-07-30 22:55:01
 # @LastEditTime: 2021-07-31 17:25:27
 # @LastEditors: lmk
 # @Description: 
### 
export GO111MODULE=auto
export GOMOD=auto
go env -w GO111MODULE=auto
go build -o mises-service ./client/cli/main.go