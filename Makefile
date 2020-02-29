PACKAGE=github.com/rebuy-de/exporter-merger
NAME=$(notdir $(PACKAGE))
BUILD_VERSION=$(shell git describe --always --dirty --tags | tr '-' '.' )
BUILD_DATE=$(shell date)
BUILD_HASH=$(shell git rev-parse HEAD)
BUILD_MACHINE=$(shell echo $$HOSTNAME)
BUILD_USER=$(shell whoami)

BUILD_FLAGS=-ldflags "\
	-X '$(PACKAGE)/cmd.BuildVersion=$(BUILD_VERSION)' \
	-X '$(PACKAGE)/cmd.BuildDate=$(BUILD_DATE)' \
	-X '$(PACKAGE)/cmd.BuildHash=$(BUILD_HASH)' \
	-X '$(PACKAGE)/cmd.BuildEnvironment=$(BUILD_USER)@$(BUILD_MACHINE)' \
"

test:
	go test ./...

build:
	go build \
		$(BUILD_FLAGS) \
		-o $(NAME)-$(BUILD_VERSION)-$(shell go env GOOS)-$(shell go env GOARCH)
	ln -sf $(NAME)-$(BUILD_VERSION)-$(shell go env GOOS)-$(shell go env GOARCH) $(NAME)

install: test
	go install \
		$(BUILD_FLAGS)

clean:
	rm -f $(NAME)*

.PHONY: build install test
