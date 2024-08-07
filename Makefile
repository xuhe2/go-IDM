# set the go project name
BINARY_NAME=go-IDM

# build the go binary
build: 
	@go build -o bin/${BINARY_NAME}

# run the go binary
run:
	@bin/${BINARY_NAME} $(ARGS)

# build and run the go binary
run-build: build run

# clean up the build artifacts
clean-bin:
	rm -f bin/${BINARY_NAME}

# clean the storage data
clean-data:
	rm -r dataDir

# test the go code
test:
	@go test ./...