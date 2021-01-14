platform = $(shell uname -s)
dest = build
tag = $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
component = ardi_$(tag)
component_path = $(dest)/$(component)
linux_objects = $(component_path)_linux_amd64 \
$(component_path)_linux_arm_5 \
$(component_path)_linux_arm_6 \
$(component_path)_linux_arm_7 \
$(component_path)_linux_arm64

darwin_object = $(component_path)_darwin_amd64

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

ifeq ($(platform),Linux)
objects := $(linux_objects)
else
objects := $(darwin_object)
endif

$(foreach o,$(objects),$(call add_zip_object,$(o)))

all: $(objects)

$(objects): $(shell find . -type f -name "*.go")
	$(eval goos=$(call get_goos, $(@)))
	$(eval goarch=$(call get_goarch, $(@)))
	$(eval goarm=$(call get_goarm, $(@)))
	GOOS=$(goos) GOARCH=$(goarch) GOARM=$(goarm) go build -ldflags '-s -w' -o $(@)

$(zips): $(objects)
	zip -j $(@) $(@:.zip=)

release: $(zips)

mock:
	go generate ./...

lint:
	golint -set_exit_status ./...

test-unit:
	mkdir -p test_artifacts
	rm -rf test_artifacts/unit.out
	go test \
		-v \
		-p 1 \
		-coverprofile test_artifacts/unit.out \
		github.com/robgonnella/ardi/v2/core \
		github.com/robgonnella/ardi/v2/util

test-int:
	mkdir -p test_artifacts
	rm -rf test_artifacts/int.out
	go test \
		-v \
		-p 1 \
		-coverprofile test_artifacts/int.out \
		github.com/robgonnella/ardi/v2/commands

test-e2e:
	./scripts/run_e2e.sh

test-all:
	rm -rf test_artifacts
	mkdir -p test_artifacts
	go test \
		-v \
		-p 1 \
		-coverprofile test_artifacts/coverage.out \
		github.com/robgonnella/ardi/v2/core \
		github.com/robgonnella/ardi/v2/util \
		github.com/robgonnella/ardi/v2/commands
	go tool cover -func test_artifacts/coverage.out
	$(MAKE) test-e2e

docs:
	go run docs/gen.go

clean:
	rm -rf build

.PHONY: all release mock docs clean
