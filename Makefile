ENV=${env}
VER=${ver}
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
	@echo run --rm $(EXEC_NAME) $(ENV)
	@./$(EXEC_NAME) --config=conf/$(ENV).yaml

# env:dev|pre|pro ver=BuildNum|TagNum
build_stage:build
	@echo  ......build stage $(ENV).$(VER) image......
	@docker rmi -f $(HarborRegistry)/$(APP):$(ENV).$(VER)
	@docker build --rm --no-cache -t $(HarborRegistry)/$(APP):$(ENV).$(VER) -f Dockerfile .
	@echo  ......build stage $(ENV).$(VER) image finish end

# env:dev|pre|pro ver=BuildNum|TagNum
push_stage:
	@echo  ......push stage $(ENV).$(VER) image......
	@docker push $(HarborRegistry)/$(APP):$(ENV).$(VER)
	@echo  ......push stage $(ENV).$(VER) image finish end

# env:dev|pre ver=build_num
deploy_stage:
	@echo ......deploy $(ENV).$(VER) ......
	@cd $(DeployPath)/overlays/$(ENV) \
		&& kustomize edit add annotation ver:$(VER) -f \
		&& kustomize edit add annotation kubesphere.io/description:'game slots version server-'$(ENV)$(VER) -f \
		&& kustomize edit set image $(HarborRegistry)/$(APP):$(ENV).$(VER) \
		&& cd - \
		&& kustomize build $(DeployPath)/overlays/$(ENV) | kubectl --kubeconfig $(config) apply -f - \
		&& kubectl --kubeconfig $(config) apply -f $(DeployPath)/ingress.yaml \

	@echo ......deploy stage $(ENV).$(VER) finish end

# env:pro ver=TagNum
deploy_pro:
	@echo  ......deploy $(ENV).$(VER) ......
	@cd $(DeployPath)/overlays/$(ENV) \
		&& kustomize edit set nameSuffix -- -acem -pro-${subst .,-,${VER}}  \
		&& kustomize edit set label ver:$(VER) \
		&& kustomize edit add annotation ver:$(VER) -f \
		&& kustomize edit add annotation kubesphere.io/description:'game slots version server-'$(ENV)$(VER) -f \
		&& kustomize edit add configmap gsv-cm --behavior=merge --from-literal ver=$(VER) \
		&& kustomize edit set image $(HarborRegistry)/$(APP):$(ENV).$(VER) \
		&& cd - \
		&& kustomize build $(DeployPath)/overlays/$(ENV) | kubectl --kubeconfig $(config) apply -f - \
		&& kubectl --kubeconfig $(config) apply -f $(DeployPath)/ingress.yaml \

	@echo  ......deploy pro.$(VER) finish end......

.PHONY: all help clean build\
		build_stage push_stage\
		deploy_stage\
		deploy_pro\

help:
	@echo "usage cmd: "
	@echo " --make clean - rm executer"
	@echo " --make build - build executer"
	@echo " --make build_stage - build stage(dev|pre|pro) image eg:env=dev|pre|pro ver=buildNum|tagNum"
	@echo " --make push_stage - push stage(dev|pre) image to harbor eg:env=dev|pre|pro ver=buildNum|tagNum"
	@echo " --make deploy_stage - deploy stage(dev|pre) ver(buildNum) eg:env=dev|pre ver=buildNum"
	@echo " --make deploy_pro - deploy stage(pro) ver(tagNum) online eg:env=pro ver=tagNum"
