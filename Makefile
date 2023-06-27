##----------------------------------------------------------------------------##
## Variables                                                                  ##
##----------------------------------------------------------------------------##

OUTPUT = ./bin/
BINARY = systags

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
	@echo "  $$ make build   - Builds systags binary for this OS"
	@echo "  $$ make clean   - Cleans and removes generated files"
	@echo "  $$ make publish - Builds artifacts for a new release"
	@echo
	@echo "DOCS"
	@echo "  Visit https://github.com/StackAdapt/systags for more"
	@echo

##----------------------------------------------------------------------------##
## Build                                                                      ##
##----------------------------------------------------------------------------##

.PHONY: build clean publish

build:
	go build -o "$(OUTPUT)$(BINARY)"

clean:
	rm -rf "$(OUTPUT)"

publish: clean
	env GOOS=linux GOARCH=amd64 go build -o "$(OUTPUT)$(BINARY)_linux_amd64/$(BINARY)"
	env GOOS=linux GOARCH=arm64 go build -o "$(OUTPUT)$(BINARY)_linux_arm64/$(BINARY)"

	tar -czf "$(OUTPUT)$(BINARY)_linux_amd64.tar.gz" -C "$(OUTPUT)" "$(BINARY)_linux_amd64"
	tar -czf "$(OUTPUT)$(BINARY)_linux_arm64.tar.gz" -C "$(OUTPUT)" "$(BINARY)_linux_arm64"
