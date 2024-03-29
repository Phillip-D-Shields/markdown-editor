BINARY_NAME=markdown.app
APP_NAME=markdown
VERSION=0.1.0

## build: build binary and package app
build:
	rm -rf ${BINARY_NAME}
	rm -f markdown-editor
	fyne package --appVersion ${VERSION} -name ${APP_NAME} -release

## run: run the app
run:
	go run .


## clean: clean up
clean:
	@echo "Cleaning up..."
	@go clean
	@rm -rf ${BINARY_NAME}
	@echo "Done."

## test: runs all tests
test:
	go test -v ./...