ENV=${env}
VSN=${vsn}
BUILD_NUM=${ver}
config=${config}
APP=gsv
EXEC_NAME=gsv
MAIN=main.go
HarborRegistry=harbor.nuclearport.com/jackpotland
DeployPath=deploy/kustomize

all: help

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

# env:dev|pre|pro
run:build
	@echo run --rm $(EXEC_NAME) $(ENV) $(VSN)
	@./$(EXEC_NAME) --config=conf/$(ENV).yaml

# env:dev|pre vsn=BuildNum
build_stage:build
	@echo  ......build stage $(ENV).$(BUILD_NUM) image......
	@docker rmi -f $(HarborRegistry)/$(APP):$(ENV).$(BUILD_NUM)
	@docker build --rm --no-cache -t $(HarborRegistry)/$(APP):$(ENV).$(BUILD_NUM) -f Dockerfile .
	@echo  ......build stage $(ENV).$(BUILD_NUM) image finish end

# env:dev|pre vsn=BuildNum
push_stage:
	@echo  ......push stage $(ENV).$(BUILD_NUM)- image......
	@docker push $(HarborRegistry)/$(APP):$(ENV).$(BUILD_NUM)
	@echo  ......push stage $(ENV).$(BUILD_NUM) image finish end

# env:dev|pre vsn=build_num
deploy_stage:
	@echo ......deploy $(ENV) vsn=$(BUILD_NUM) ......;source /etc/profile ;which aws ;echo $PATH
	@cd $(DeployPath)/overlays/$(ENV) \
		&& kustomize edit set namesuffix -- -$(ENV)-v${BUILD_NUM} \
		&& kustomize edit set label vsn:$(BUILD_NUM) \
		&& kustomize edit set annotation vsn:$(BUILD_NUM)\
		&& kustomize edit add annotation kubesphere.io/description:'game slots version server '$(ENV)-$(VSN) \
		&& kustomize edit set image $(HarborRegistry)/$(APP):$(ENV).$(VSN) \
		&& kustomize edit add configmap gsv-cm --behavior=merge --from-literal vsn=$(VSN) \
		&& cd - \
		&& kustomize build $(DeployPath)/overlays/$(ENV) | kubectl --kubeconfig $(config) apply -f - \
		&& kubectl --kubeconfig $(config) apply -f $(DeployPath)/ingress.yaml \

	@echo  ......deploy $(ENV) vsn=$(BUILD_NUM) finish end......

# vsn:1.0.0
build_pro:build
	@echo  ......build stage pro vsn=$(VSN) image......
	@docker rmi -f $(HarborRegistry)/$(APP):pro.$(VSN)
	@docker build --rm --no-cache -t $(HarborRegistry)/$(APP):pro.$(VSN) -f Dockerfile .
	@echo  ......build stage pro vsn=$(VSN) image finish end

# vsn:1.0.0
push_pro:
	@echo ......push stage pro vsn=$(VSN) image......
	@docker push $(HarborRegistry)/$(APP):pro.$(VSN)
	@echo ......push stage pro vsn=$(VSN)  image finish end

# vsn:1.0.0
deploy_pro:
	@echo  ......deploy pro vsn=$(VSN) ......;source /etc/profile ;which aws ;echo $PATH
	@cd $(DeployPath)/overlays/pro \
		&& kustomize edit set namesuffix -- -pro-v${subst .,-,${VSN}} \
		&& kustomize edit set label vsn:$(VSN) \
		&& kustomize edit set annotation vsn:$(VSN)\
		&& kustomize edit add annotation kubesphere.io/description:'game slots version server pro-'$(VSN) \
		&& kustomize edit set image $(HarborRegistry)/$(APP):pro.$(VSN) \
		&& kustomize edit add configmap gsv-cm --behavior=merge --from-literal vsn=$(VSN) \
		&& cd - \
		&& kustomize build $(DeployPath)/overlays/pro | kubectl --kubeconfig $(config) apply -f - \
		&& kubectl --kubeconfig $(config) apply -f $(DeployPath)/ingress.yaml \

	@echo  ......deploy pro vsn=$(VSN) finish end......

.PHONY: all help clean build\
		build_stage push_stage deploy_stage\
		build_pro push_pro deploy_pro\

help:
	@echo "usage cmd: "
	@echo " --make clean - rm executer"
	@echo " --make build - build executer"
	@echo " --make build_stage - build stage(dev|pre) image eg:env=dev vsn=BuildNum|env=pre vsn=BuildNum"
	@echo " --make push_stage - push stage(dev|pre) image to harbor eg:env=dev vsn=BuildNum|env=pre vsn=BuildNum"
	@echo " --make deploy_stage - deploy stage(dev|pre) eg:env=dev vsn=BuildNum|env=pre vsn=BuildNum"
	@echo " --make build_pro - build pro image eg:vsn=1.0.0"
	@echo " --make push_pro - push pro image to harbor eg:vsn=1.0.0"
	@echo " --make deploy_pro - deploy pro vsn=1.0.0 online eg:vsn=1.0.0"
