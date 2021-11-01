VSN=${vsn}
PROJECT=${project}
EXE_NAME=gs_vsn
MAIN=main.go
HarborRegistry=harbor.nuclearport.com/jackpotland/${PROJECT}/

.PHONY: all help fmt clean build run build_image push_image

all: help

fmt:
	@echo fmt......
	@go fmt ./
	@go vet ./
	@echo fmt end

clean:
	@echo clean......
	@if [ -f ${EXE_NAME} ] ; then rm ${EXE_NAME} ; fi
	@echo clean finish end

build:clean fmt
	@echo build......
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	@go build -a -installsuffix cgo -o $(EXE_NAME) $(MAIN)
	@echo build finish end

run:build
	@echo run ${EXE_NAME}
	@./$(EXE_NAME) --config=conf/dev.yaml

build_image:build
	@echo build image......
	@docker build -no-cache -t $(HarborRegistry)$(PROJECT):$(VSN) -f Dockerfile .
	@docker tag $(PROJECT):$(VSN) $(HarborRegistry)$(PROJECT):$(VSN)
	@echo build image finish end

push_image: build_image
	@echo push image......
	@docker push $(HarborRegistry)$(PROJECT):$(VSN)
	@echo push image finish end

help:
	@echo "usage cmd: "
	@echo " --make build - generate exe binary file"
	@echo " --make run - run main"
	@echo " --make clean - rm exe binary file"
	@echo " --make gotool - fmt vet"
