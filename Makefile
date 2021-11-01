EXE_NAME=gs_vsn
MAIN=main.go

.PHONY: all help fmt clean build run

all: help

fmt:
	@echo fmt
	@go fmt ./
	@go vet ./

clean:
	@echo clean
	@if [ -f ${EXE_NAME} ] ; then rm ${EXE_NAME} ; fi
	@echo clean finish

build:clean fmt
	@echo build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	@go build -a -installsuffix cgo -o $(EXE_NAME) $(MAIN)
	@echo build finish

run:build
	@echo run ${EXE_NAME}
	@./$(EXE_NAME) --config=conf/dev.yaml

help:
	@echo "usage cmd: "
	@echo " --make build - generate exe binary file"
	@echo " --make run - run main"
	@echo " --make clean - rm exe exe binary file"
	@echo " --make gotool - fmt vet"
