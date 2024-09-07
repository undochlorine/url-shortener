# Build targets
.PHONY: proto

all: proto

proto:
	rm -f ./**/*.pb.go
	protoc --proto_path=./ --go_out=./ --go-grpc_out=./ --go_opt=paths=source_relative */*.proto