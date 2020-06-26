dest = build
tag = $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
component = ardi_$(tag)
component_path = $(dest)/$(component)
objects = $(component_path)_darwin_amd64 \
$(component_path)_linux_amd64 \
$(component_path)_linux_arm_5 \
$(component_path)_linux_arm_6 \
$(component_path)_linux_arm_7 \
$(component_path)_linux_arm64

define get_goos
$(word 3, $(subst _, ,$1))
endef

define get_goarch
$(word 4, $(subst _, ,$1))
endef

define get_goarm
$(word 5, $(subst _, ,$1))
endef

all: $(objects)

$(objects): $(shell find . -type f -name "*.go")
	$(eval goos=$(call get_goos, $(@)))
	$(eval goarch=$(call get_goarch, $(@)))
	$(eval goarm=$(call get_goarm, $(@)))
	GOOS=$(goos) GOARCH=$(goarch) GOARM=$(goarm) go build -ldflags '-s -w' -o $(@)

mock:
	go generate ./...

test:
	go test -v ./...

docs:
	go run docs/gen.go

clean:
	rm -rf build

.PHONY: all mock docs clean
