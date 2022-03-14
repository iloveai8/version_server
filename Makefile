ENV=${env}
VER=${ver}
VSN=${vsn}
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

# env:dev|pre ver=BuildNum
build_stage:build
	@echo  ......build stage $(ENV).$(VER) image......
	@docker rmi -f $(HarborRegistry)/$(APP):$(ENV).$(VER)
	@docker build --rm --no-cache -t $(HarborRegistry)/$(APP):$(ENV).$(VER) -f Dockerfile .
	@echo  ......build stage $(ENV).$(VER) image finish end

# env:dev|pre ver=BuildNum
push_stage:
	@echo  ......push stage $(ENV).$(VER) image......
	@docker push $(HarborRegistry)/$(APP):$(ENV).$(VER)
	@echo  ......push stage $(ENV).$(VER) image finish end

# env:dev|pre ver=build_num
deploy_stage:
	@echo ......deploy $(ENV) $(VER) ......
	@cd $(DeployPath)/overlays/$(ENV) \
		&& kustomize edit add annotation ver:$(VER) \
		&& kustomize edit add annotation kubesphere.io/description:'game slots version server' $(ENV) '-' $(VER) \
		&& kustomize edit set image $(HarborRegistry)/$(APP):$(ENV).$(VER) \
		&& cd - \
		&& kustomize build $(DeployPath)/overlays/$(ENV) | kubectl --kubeconfig $(config) apply -f - \
		&& kubectl --kubeconfig $(config) apply -f $(DeployPath)/ingress.yaml \

	@echo ......deploy stage $(ENV) $(VER) finish end

# vsn:1.0.0
build_pro:build
	@echo  ......build stage pro $(VER) image......
	@docker rmi -f $(HarborRegistry)/$(APP):pro.$(VER)
	@docker build --rm --no-cache -t $(HarborRegistry)/$(APP):pro.$(VER) -f Dockerfile .
	@echo  ......build stage pro $(VER) image finish end

# vsn:1.0.0
push_pro:
	@echo ......push stage pro $(VER) image......
	@docker push $(HarborRegistry)/$(APP):pro.$(VSN)
	@echo ......push stage pro $(VER) image finish end

# vsn:1.0.0
deploy_pro:
	@echo  ......deploy pro $(VER) ......;source /etc/profile ;which aws ;echo $PATH
	@cd $(DeployPath)/overlays/pro \
		&& kustomize edit set namesuffix -- -pro-v${subst .,-,${VER}} \
		&& kustomize edit set label ver:$(VER) \
		&& kustomize edit set add ver:$(VER)\
		&& kustomize edit add annotation kubesphere.io/description:'game slots version server pro-'$(VER) \
		&& kustomize edit set image $(HarborRegistry)/$(APP):pro.$(VER) \
		&& kustomize edit add configmap gsv-cm --behavior=merge --from-literal ver='$(VER)' \
		&& cd - \
		&& kustomize build $(DeployPath)/overlays/pro | kubectl --kubeconfig $(config) apply -f - \
		&& kubectl --kubeconfig $(config) apply -f $(DeployPath)/ingress.yaml \

	@echo  ......deploy pro $(VER) finish end......

.PHONY: all help clean build\
		build_stage push_stage deploy_stage\
		build_pro push_pro deploy_pro\

help:
	@echo "usage cmd: "
	@echo " --make clean - rm executer"
	@echo " --make build - build executer"
	@echo " --make build_stage - build stage(dev|pre) image eg:env=dev ver=BuildNum|env=pre ver=BuildNum"
	@echo " --make push_stage - push stage(dev|pre) image to harbor eg:env=dev ver=BuildNum|env=pre ver=BuildNum"
	@echo " --make deploy_stage - deploy stage(dev|pre) eg:env=dev ver=BuildNum|env=pre ver=BuildNum"
	@echo " --make build_pro - build pro image eg:vsn=1.0.0"
	@echo " --make push_pro - push pro image to harbor eg:vsn=1.0.0"
	@echo " --make deploy_pro - deploy pro vsn=1.0.0 online eg:vsn=1.0.0"
