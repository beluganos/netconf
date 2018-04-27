
SUBDIRS = vty gobgp sysrepo openconfig lxd sysctl netplan property signal net xml

.PHONY: go-test

go-test:
	go test -coverprofile=cover.out

check-local: go-test
