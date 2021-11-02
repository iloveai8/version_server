ENV=dev
VSN=1.0.0
TEAM=jackpotland
PROJECT=game_slots_vsn
EXEC_NAME=game_slots_vsn
MAIN=main.go
HarborRegistry=harbor.nuclearport.com

.PHONY: all help fmt clean build run build_image push_image

all: help
	@echo "vsn is $(VSN)"
	@echo "env is $(ENV)"

fmt:
	@go fmt ./
	@go vet ./

clean:
	@if [ -f $(EXEC_NAME) ] ; then rm $(EXEC_NAME) ; fi

build:clean fmt
	@echo build......
	@set CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	@go build -a -installsuffix cgo -o $(EXEC_NAME) $(MAIN)
	@echo build finish end

run:build
	@echo run $(EXEC_NAME) $(ENV) $(VSN)
	@./$(EXEC_NAME) --config=conf/$(ENV).yaml

build_image:build
	@echo build image......
	@docker build --no-cache -t $(TEAM)$(PROJECT):$(VSN) -f Dockerfile .
	@echo build image finish end

push_image: build_image
	@echo push image......
	@docker push $(HarborRegistry)$(TEAM)$(PROJECT):$(VSN)
	@echo push image finish end

help:
	@echo "usage cmd: "
	@echo " --make fmt - fmt vet"
	@echo " --make clean - rm exe binary file"
	@echo " --make build - generate exe binary file"
	@echo " --make run - build and run main"
	@echo " --make build_image - build docker image"
	@echo " --make push_image - build docker image and push harbor"
