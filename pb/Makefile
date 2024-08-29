# Build targets
.PHONY: proto

all: proto mocks

proto:
	rm -f ./**/*.pb.go
	protoc --proto_path=./ --go_out=./ --go-grpc_out=./ --go_opt=paths=source_relative */*.proto
	protoc-go-inject-tag -input=*/*.pb.go

projects = $(shell dirname */*.proto)

mocks:
	rm -f pb/*.go
	$(foreach project, $(projects), \
		test -e ./$(project)/$(project)_grpc.pb.go && \
		mockgen --package=$(project) \
			-source=./$(project)/$(project)_grpc.pb.go \
			-destination=./$(project)/$(project)_mock.go ;)


