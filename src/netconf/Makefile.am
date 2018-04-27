# -*- coding: utf-8 -*-

SUBDIRS = lib app

PACKAGES = netconf/...

go-fmt:
	goimports -w .

go-vet:
	go vet ${PACKAGES}

go-build:
	go build ${PACKAGES}

go-install:
	go install ${PACKAGES}

all-local: go-fmt go-vet go-build

install-exec-local: go-install
