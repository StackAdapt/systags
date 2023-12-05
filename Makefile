##----------------------------------------------------------------------------##
## Variables                                                                  ##
##----------------------------------------------------------------------------##

OUTPUT  = ./bin/
INSTALL = /usr/local/bin/
BINARY  = systags



##----------------------------------------------------------------------------##
## Help                                                                       ##
##----------------------------------------------------------------------------##

.PHONY: help

help:
	@echo
	@echo "WELCOME TO SYSTAGS"
	@echo "------------------"
	@echo
	@echo "MAKE"
	@echo "  $$ make help    - Prints out these help instructions"
	@echo "  $$ make build   - Builds main binary in release mode"
	@echo "  $$ make debug   - Builds main binary in debug mode"
	@echo "  $$ make clean   - Cleans and removes generated files"
	@echo "  $$ make install - Builds and installs binary on system"
	@echo "  $$ make remove  - Removes installed binary from system"
	@echo "  $$ make publish - Builds artifacts for a new release"
	@echo
	@echo "DOCS"
	@echo "  Visit https://github.com/StackAdapt/systags for more"
	@echo



##----------------------------------------------------------------------------##
## Build                                                                      ##
##----------------------------------------------------------------------------##

.PHONY: build debug clean

build:
	go build -ldflags "-s -w" -o "$(OUTPUT)$(BINARY)"

debug:
	# Include -gcflags to improve the experience with GDB
	# particularly when needing to print variable values.
	# Also try to detect race conditions using -race flag.
	go build -o "$(OUTPUT)$(BINARY)" -gcflags="all=-N -l" -race

clean:
	rm -rf "$(OUTPUT)"



##----------------------------------------------------------------------------##
## Install                                                                    ##
##----------------------------------------------------------------------------##

.PHONY: install remove

install: build
	mkdir -p "$(INSTALL)" # In case directory is missing
	install -p -m 0755 "$(OUTPUT)$(BINARY)" "$(INSTALL)"

remove:
	# Remove installed binary
	rm -f "$(INSTALL)$(BINARY)"



##----------------------------------------------------------------------------##
## Publish                                                                    ##
##----------------------------------------------------------------------------##

.PHONY: _publish_linux_amd64 _publish_linux_arm64 publish

_publish_linux_amd64: clean
	env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "$(OUTPUT)$(BINARY)_linux_amd64/$(BINARY)"
	tar -czf "$(OUTPUT)$(BINARY)_linux_amd64.tar.gz" -C "$(OUTPUT)" "$(BINARY)_linux_amd64"

_publish_linux_arm64: clean
	env GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o "$(OUTPUT)$(BINARY)_linux_arm64/$(BINARY)"
	tar -czf "$(OUTPUT)$(BINARY)_linux_arm64.tar.gz" -C "$(OUTPUT)" "$(BINARY)_linux_arm64"

publish: _publish_linux_amd64 _publish_linux_arm64
