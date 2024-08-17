proto:
	export GOROOT=/usr/local/go
	export GOPATH=$HOME/go
	export GOBIN=$GOPATH/bin
	export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
	protoc --go_out=./proto --go-grpc_out=./proto --go-grpc_opt=require_unimplemented_servers=false ./proto/*.proto
	# protoc --go_out=./proto --go-grpc_out=./proto --go-grpc_opt=require_unimplemented_servers=false ./proto/auth.proto
proto:
	protoc --go_out=./proto --go-grpc_out=./proto --go-grpc_opt=require_unimplemented_servers=false ./proto/*.proto
server:
	go run userproto/server/server.go
client:
	go run userproto/client/user_client.go
run:
	go run cmd/main.go

task:
	task :start -w --interval=500ms
	
.PHONY: proto server client task
