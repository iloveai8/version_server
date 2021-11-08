ENV=${env}
VSN=${vsn}
PROJECT=gsv
EXEC_NAME=gsv
MAIN=main.go
HarborRegistry=harbor.nuclearport.com/jackpotland
DeployPath=deploy/kustomize/overlays

.PHONY: all help fmt clean build run build_image push_image deploy

all: help
	@echo "vsn is $(VSN)"
	@echo "env is $(ENV)"jackpotland

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
	@docker build --no-cache -t $(HarborRegistry)/$(PROJECT):$(VSN) -f Dockerfile .
	@echo build image finish end

push_image: build_image
	@echo push image......
	@docker push $(HarborRegistry)/$(PROJECT):$(VSN)
	@echo push image finish end

# env:dev|pre|pro vsn:1.0.0
deploy:
	@echo deploy env:$(ENV) vsn:$(VSN) ......
	@echo $(DeployPath)/$(ENV)
	@cd $(DeployPath)/$(ENV) \
		&& pwd \
		&& kustomize edit set namesuffix -- -$(ENV)-$(VSN) \
		&& kustomize edit set label app:$(PROJECT)-$(ENV)-$(VSN) \
		&& kustomize edit set image $(HarborRegistry)/$(PROJECT):$(VSN) \
		&& kustomize edit add configmap gsv-config-$(ENV)-$(VSN) --from-literal=vsn=$(VSN) \
		&& kustomize build $(DeployPath)/$(ENV) \
		#| kubectl --kubeconfig $(config) apply -f -
	@echo deploy env:$(ENV) vsn:$(VSN) finish end

help:
	@echo "usage cmd: "
	@echo " --make fmt - fmt vet"
	@echo " --make clean - rm exe binary file"
	@echo " --make build - generate exe binary file"
	@echo " --make run - build and run main"
	@echo " --make build_image - build docker image"
	@echo " --make push_image - build docker image and push harbor"
	@echo " --deploy server env and vsn to online"
