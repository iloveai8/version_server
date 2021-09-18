CONFIG=./conf/dev.yaml
BINARY="GSVsnServer"
MAIN=main.go

.PHONY: all build run gotool clean help
all: help
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $(BINARY) $(MAIN)
run:
	@go run ./ --config=${CONFIG}
gotool:
	go fmt ./
	go vet ./
clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
help:
	@echo "usage cmd: "
	@echo " --make build - generate exe binary file"
	@echo " --make run - run main"
	@echo " --make clean - rm exe exe binary file"
	@echo " --make gotool - fmt vet"
