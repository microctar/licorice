.DEFAULT_GOAL=current

NAME=licorice
BINDIR=bin
VERSION=$(shell git describe --tags || echo "unknown version")
GITCOMMIT=$(shell git rev-parse --short HEAD || echo "unsupported")
BUILDTIME=$(shell date -u)

#  -trimpath => remove all file system paths from executable
#  -ldflags -X => set package variable

GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-X "github.com/microctar/licorice/app/constant.Version=$(VERSION)" \
				-X "github.com/microctar/licorice/app/constant.BuildTime=$(BUILDTIME)" \
				-X "github.com/microctar/licorice/app/constant.GitCommit=$(GITCOMMIT)" \
				-w -s -buildid='

PLATFORM_LIST = \
		linux-386 \
		linux-amd64 \
		linux-amd64-v3 \
		linux-armv5 \
		linux-armv6 \
		linux-armv7 \
		linux-armv8 \
		freebsd-386 \
		freebsd-amd64 \
		freebsd-amd64-v3 \
		freebsd-arm64

all: linux-amd64 freebsd-amd64

current: $(shell uname -sm  | tr [A-Z'\040'] [a-z'\055'])

linux-i386: linux-386

linux-i686: linux-386

linux-386:
	GOARCH=386 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-x86_64: linux-amd64

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-amd64-v3:
	GOARCH=amd64 GOOS=linux GOAMD64=v3 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-armv5:
	GOARCH=arm GOOS=linux GOARM=5 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-armv6:
	GOARCH=arm GOOS=linux GOARM=6 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-armv7:
	GOARCH=arm GOOS=linux GOARM=6 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-aarch64: linux-armv8

linux-armv8:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

freebsd-i386: freebsd-386

freebsd-386:
	GOARCH=386 GOOS=freebsd $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

freebsd-amd64:
	GOARCH=amd64 GOOS=freebsd $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

freebsd-aarch64: freebsd-arm64

freebsd-arm64:
	GOARCH=arm64 GOOS=freebsd $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

gz_releases=$(addsuffix .gz, $(PLATFORM_LIST))

all-arch: $(PLATFORM_LIST)

lint:
	GOOS=linux golangci-lint run ./...
	GOOS=freebsd golangci-lint run ./...

clean:
	rm $(BINDIR)/*
