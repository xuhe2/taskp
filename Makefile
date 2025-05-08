ServerDir = ./cmd/server/
ClientDir = ./cmd/client/

.PHONY: build
build:
	go build -o ./bin/server ${ServerDir}
	go build -o ./bin/client ${ClientDir}

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: gen-grpc
gen-grpc:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	./netapi/task.proto

