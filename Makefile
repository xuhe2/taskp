ServerFile = ./cmd/server/main.go
ClientFile = ./cmd/client/main.go

.PHONY: build
build:
	go build -o ./bin/server ${ServerFile}
	go build -o ./bin/client ${ClientFile}

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: gen-grpc
gen-grpc:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	./netapi/task.proto

