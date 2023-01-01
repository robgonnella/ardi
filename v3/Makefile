platform ?= $(shell uname -s)
dest = build
tag = $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
ci ?= false

#### Build Objects ####
component = ardi_$(tag)
component_path = $(dest)/$(component)
linux_objects = $(component_path)_linux_amd64 \
$(component_path)_linux_arm_5 \
$(component_path)_linux_arm_6 \
$(component_path)_linux_arm_7 \
$(component_path)_linux_arm64

darwin_object = $(component_path)_darwin_amd64

#### Zip File Objects ####

# see gather objects section below
zips =

#### Gather Objects ####

ifeq ($(platform),Linux)
objects := $(linux_objects)
else
objects := $(darwin_object)
endif

$(foreach o,$(objects), $(eval zips += $(o).zip))

#### Test Objects ####
test_output_dir = test_artifacts
coverage_out = $(test_output_dir)/coverage.txt

#### Helper Functions ####

define get_goos
$(word 3, $(subst _, ,$1))
endef

define get_goarch
$(word 4, $(subst _, ,$1))
endef

define get_goarm
$(word 5, $(subst _, ,$1))
endef

#### Rules Section ####

.PHONY: all
all: $(objects)

$(objects): $(shell find . -type f -name "*.go")
	$(eval goos=$(call get_goos, $(@)))
	$(eval goarch=$(call get_goarch, $(@)))
	$(eval goarm=$(call get_goarm, $(@)))
	GOOS=$(goos) GOARCH=$(goarch) GOARM=$(goarm) go build -ldflags '-s -w' -o $(@)

$(zips): $(objects)
	zip -j $(@) $(@:.zip=)

.PHONY: release
release: $(zips)

.PHONY: mock
mock:
	go generate ./...

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY:
go-test:
	mkdir -p $(test_output_dir)
	rm -rf $(coverage_out)
	go test \
		-v \
		-p 1 \
		-coverprofile $(coverage_out) \
		-covermode=atomic \
		./...

.PHONY: test-e2e
test-e2e:
	./scripts/run_e2e.sh

.PHONY: test-all
test-all: lint
	$(MAKE) go-test
	$(MAKE) test-e2e
	$(MAKE) coverage

.PHONY: coverage
coverage:
	go tool cover -func $(coverage_out)

.PHONY: test-report
test-report:
	go tool cover -html=$(coverage_out)

.PHONY: docs
docs:
	go run docs/gen.go

.PHONY: deps
deps:
	go install github.com/golang/mock/mockgen@latest
	go install golang.org/x/lint/golint@latest


.PHONY: clean
clean:
	rm -rf build
