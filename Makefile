PROTO_DIR=./proto
GO_OUT_DIR=./task-service/types
RUST_OUT_DIR=./notification-service/src/types

.PHONY: proto-gen-go proto-gen-rust

proto-gen-go:
	rm -rf $(GO_OUT_DIR)/*
	protoc -I=$(PROTO_DIR) --go_out=$(GO_OUT_DIR) $(PROTO_DIR)/*
	mv $(GO_OUT_DIR)/github.com/stringintech/task-broker/types/* $(GO_OUT_DIR)/
	rm -rf $(GO_OUT_DIR)/github.com

proto-gen-rust:
	cd notification-service && cargo clean && cargo build #TODO? should not have to clean to detect proto changes
