PROTO = service.proto

proto:
	protoc --go_out=. $(PROTO)
	protoc --nats-rpc_out=subject=$(NAME):. $(PROTO)
	protoc --nats-rpc-cli_out=subject=$(NAME):cmd/cli $(PROTO)
