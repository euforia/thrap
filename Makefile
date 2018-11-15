
NAME = thrap
BUILDTIME = $(shell date +%Y-%m-%dT%T%z)

# These options are to allow for cross-platform support
DIST_OPTS = -a -tags netgo -installsuffix netgo
# This sets the version and build time in the binary
LD_OPTS = -ldflags="-X main._version=$(STACK_VERSION) -X main._buildtime=$(BUILDTIME) -w"
# CGO_ENABLED is set to be zero to cross-platform support
BUILD_CMD = CGO_ENABLED=0 go build $(LD_OPTS)

SOURCE_FILES = $(shell ls ./cmd/*.go | grep -v _test.go)
SOURCE_PACKAGES = $(shell go list ./... | grep -v /vendor/ | grep -v /crt)

clean:
	rm -rf ./dist
	rm -rf ./ui/build
	rm -f ./pkg/api/ui.go
	rm -f ./pkg/api/swagger.go

.PHONY: deps
deps:
	go get github.com/go-swagger/go-swagger/cmd/swagger
	go get github.com/jteeuwen/go-bindata
	go get github.com/golang/dep/cmd/dep

.PHONY: test
test:
	go test -cover $(SOURCE_PACKAGES)
	
# Generate swagger docs
ui/build/swagger.json:
	swagger generate spec -i swagger.yml -o ./ui/build/swagger.json

# Build react app
ui/build:
	cd ./ui/ && yarn --verbose build

# Bindata the ui
bindata:
	cd ./ui/ && go-bindata -pkg api -o ../pkg/api/ui.go build/...

# Build complete ui
ui: ui/build ui/build/swagger.json bindata

$(NAME):
	$(BUILD_CMD) -o $(NAME) $(SOURCE_FILES)

dist/$(NAME)-linux:
	GOOS=linux $(BUILD_CMD) $(DIST_OPTS) -o dist/$(NAME)-linux $(SOURCE_FILES)

# linux is always done last as it is used in the container build
dist:
	mkdir dist
	for goos in darwin linux; do \
	GOOS=$${goos} $(BUILD_CMD) $(DIST_OPTS) -o dist/$(NAME) $(SOURCE_FILES); \
	tar -czf dist/$(NAME)-$${goos}.tgz -C dist $(NAME); done

	
	
