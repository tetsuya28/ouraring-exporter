GO := go
BIN := $(abspath ./bin)
PATH := $(abspath ./bin):$(PATH)

$(BIN)/oapi-codegen:
	GOBIN=$(BIN) $(GO) install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

.PHONY: generate
generate: $(BIN)/oapi-codegen
generate:
	oapi-codegen -package ouraring https://cloud.ouraring.com/v2/static/json/openapi-1.11.json > ouraring/ouraring_gen.go
