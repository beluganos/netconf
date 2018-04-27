# -*- coding: utf-8 -*-

# SUBDIRS =

PROTOS = rpcapi.proto

.PHONY: proto go-test

go-test:
	go test -coverprofile=cover.out

proto-golang: $(PROTOS)
	protoc -I@FF_HOME@/src -I. --go_out=plugins=grpc:. $<

proto-python: $(PROTOS)
	python -m grpc_tools.protoc -I@FF_HOME@/src -I. --python_out=. --grpc_python_out=. $<

proto: proto-golang

all-local: proto

check-local: go-test

clean-local:
	-rm $(GO_SRCS)
