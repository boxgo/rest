#
# Telescreen Makefile
#


GO          ?=  go
BASEDIR      =  $(shell pwd)
GOOS         =  $(shell go env GOOS)
GOARCH       =  $(shell go env GOARCH)
GOPATH       =  $(shell go env GOPATH)
GOFILE       =  applet
pkgs         =  $(shell $(GO) list ./... | grep -v /vendor/)

export CONFIG_ROOTPATH=$(BASEDIR)

all: clean format vet


# go format
format:
	@echo ">> formatting code"
	$(GO) fmt $(pkgs)


# go vet
vet:
	@echo ">> vetting code"
	$(GO) vet $(pkgs)


# go clean
clean:
	@echo ">> cleaning project."
	rm -f ${BASEDIR}/${GOFILE}
	$(GO) clean -i


# go doc
doc:
	@echo ">> generate api doc"
	swag init -g src/controllers


# go run dev
dev:
	cd example && gowatch -o .debug
