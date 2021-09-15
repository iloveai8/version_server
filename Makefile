config=./conf/dev.yaml
exeName=GSVsnServer
exeFunc=main.go

all: help

###===================================================================
### run
###===================================================================
.PHONY: run
run:
	@echo "run config:${config}"
	go run -race main.go --config=${config}

###===================================================================
### build
###===================================================================
.PHONY: build
build:
	@echo "build start...."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(exeName) $(exeFunc)
	@echo "build end! generate execute file size:`du -sh $(exeName)`"

###===================================================================
### other
###===================================================================
.PHONY: help
help:
	 echo "Usage make run|build"