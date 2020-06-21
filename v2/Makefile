version = $(shell git rev-parse --short HEAD)
all: build_osx build_linux build_arm5 build_arm6 build_arm7 build_arm8

mock:
	go generate ./...

build:
ifdef arm
	GOOS=$(os) GOARCH=$(arch) GOARM=$(arm) go build -ldflags '-s -w' -o build/ardi_$(version)_$(os)_$(arch)$(arm)
else
	GOOS=$(os) GOARCH=$(arch) go build -ldflags '-s -w' -o build/ardi_$(version)_$(os)_$(arch)
endif

build_osx:
	os=darwin arch=amd64 $(MAKE) build

build_linux:
	os=linux arch=amd64 $(MAKE) build

build_arm5:
	os=linux arch=arm arm=5 $(MAKE) build

build_arm6:
	os=linux arch=arm arm=6 $(MAKE) build

build_arm7:
	os=linux arch=arm arm=7 $(MAKE) build

build_arm8:
	os=linux arch=arm64 $(MAKE) build

.PHONY: all \
build \
build_osx \
build_linux \
build_arm5 \
build_arm6 \
build_arm7 \
build_arm8
