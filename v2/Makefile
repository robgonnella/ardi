include_darwin ?= false
dest = build
tag = $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
component = ardi_$(tag)
component_path = $(dest)/$(component)
objects = $(component_path)_linux_amd64 \
$(component_path)_linux_arm_5 \
$(component_path)_linux_arm_6 \
$(component_path)_linux_arm_7 \
$(component_path)_linux_arm64

darwin = $(component_path)_darwin_amd64

ifneq ($(include_darwin),false)
objects += $(darwin)
endif

define add_zip_object
$(eval zips += $(1).zip)
endef

define get_goos
$(word 3, $(subst _, ,$1))
endef

define get_goarch
$(word 4, $(subst _, ,$1))
endef

define get_goarm
$(word 5, $(subst _, ,$1))
endef

$(foreach o,$(objects),$(call add_zip_object,$(o)))

all: $(objects)

$(objects): $(shell find . -type f -name "*.go")
	$(eval goos=$(call get_goos, $(@)))
	$(eval goarch=$(call get_goarch, $(@)))
	$(eval goarm=$(call get_goarm, $(@)))
	GOOS=$(goos) GOARCH=$(goarch) GOARM=$(goarm) go build -ldflags '-s -w' -o $(@)

mock:
	go generate ./...

lint:
	golint -set_exit_status ./...

test-unit:
	mkdir test_artifacts || true
	rm -rf test_artifacts/unit.out || true
	go test \
		-v \
		-p 1 \
		-coverprofile test_artifacts/unit.out \
		github.com/robgonnella/ardi/v2/core \
		github.com/robgonnella/ardi/v2/util

test-int:
	mkdir test_artifacts || true
	rm -rf test_artifacts/int.out || true
	go test -v -coverprofile test_artifacts/int.out -p 1 github.com/robgonnella/ardi/v2/commands

test-e2e:
	./scripts/run_e2e.sh

test-all:
	$(MAKE) test-unit
	$(MAKE) test-int
	$(MAKE) test-e2e
	./scripts/coverage.sh

docs:
	go run docs/gen.go

clean:
	rm -rf build

.PHONY: all release mock docs clean
