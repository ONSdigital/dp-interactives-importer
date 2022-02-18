USER=$(shell id -u)
GROUP=$(shell id -g)
GOPATH=$(shell go env GOPATH)
PWD=$(shell pwd)
YAML_SPEC=$(HOME)/workspace/methods/ons/dp-upload-service/swagger.yaml

.PHONY: go-swagger
go-swagger:
	docker run --rm --user $(USER):$(GROUP) -e GOPATH=$(GOPATH):/go -v $(HOME):$(HOME) -w $(PWD) \
		quay.io/goswagger/swagger generate client -f $(YAML_SPEC) -A dp-upload-service

.PHONY: openapi
openapi:
	docker run --rm -v $(HOME):$(HOME) -v "$(PWD):/local" openapitools/openapi-generator-cli generate \
      -i $(YAML_SPEC) \
      -g go \
      -o /local/go; \
    mv $(PWD)/go/go.mod $(PWD)/go/go.mod.bak; \
    rm -f $(PWD)/go/go.sum