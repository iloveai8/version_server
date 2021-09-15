BINARY_NAME=GSVsnServer
MAIN_PATH=main.go

all: help

.PHONY:run
run:
	go run -race main.go --config=./conf/dev.yaml

.PHONY:build
build:
	@echo build
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "execute file size:`du -sh $(BINARY_NAME)`"

###===================================================================
### other
###===================================================================
.PHONY: help
help:
	 echo "Usage make run|build"