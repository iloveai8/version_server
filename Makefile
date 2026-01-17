# Makefile for game_slots_vsn
#
# 使用方法：
#   make build        - 编译程序
#   make run          - 运行程序（开发环境）
#   make test         - 运行所有测试
#   make coverage     - 生成测试覆盖率报告
#   make clean        - 清理编译产物
#   make fmt          - 格式化代码
#   make vet          - 运行go vet检查

# ==================== 变量定义 ====================

# 应用名称
APP_NAME := game_slots_vsn
EXEC_NAME := gsv

# 版本号
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v2.0.0")

# 构建时间
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Git提交哈希
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go编译参数
GO := go
GOFLAGS := -v

# 构建目录
BUILD_DIR := build
BIN_DIR := bin

# Docker配置
HarborRegistry := harbor.nuclearport.com/jackpotland
DeployPath := deploy/kustomize

# ==================== 目标 ====================

.PHONY: all
all: build

## build: 编译程序
.PHONY: build
build:
	@echo "编译 $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(EXEC_NAME) ./main.go
	@echo "编译完成: $(BIN_DIR)/$(EXEC_NAME)"

## build-linux: 编译Linux版本
.PHONY: build-linux
build-linux:
	@echo "编译Linux版本..."
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(EXEC_NAME) ./main.go
	@echo "编译完成: $(BIN_DIR)/$(EXEC_NAME)"

## run: 运行程序（开发环境）
.PHONY: run
run: build
	@echo "运行开发环境..."
	APP_ENV=dev $(BIN_DIR)/$(EXEC_NAME) run

## run-prod: 运行程序（生产环境）
.PHONY: run-prod
run-prod: build
	@echo "运行生产环境..."
	APP_ENV=pro $(BIN_DIR)/$(EXEC_NAME) run

## test: 运行所有测试
.PHONY: test
test:
	@echo "运行测试..."
	$(GO) test -v -race -count=1 ./...

## test-coverage: 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "运行测试并生成覆盖率..."
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

## coverage: 生成覆盖率报告（HTML）
.PHONY: coverage
coverage:
	@echo "生成覆盖率报告..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告: coverage.html"

## coverage-summary: 查看覆盖率摘要
.PHONY: coverage-summary
coverage-summary:
	@echo "生成覆盖率摘要..."
	$(GO) test ./... -cover | grep -E "^ok|FAIL"

## clean: 清理编译产物
.PHONY: clean
clean:
	@echo "清理编译产物..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(BUILD_DIR)
	@rm -f $(EXEC_NAME)
	@rm -f coverage.out coverage.html
	@echo "清理完成"

## fmt: 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	$(GO) fmt ./...
	$(GO) vet ./...
	@echo "格式化完成"

## vet: 运行go vet检查
.PHONY: vet
vet:
	@echo "运行go vet..."
	$(GO) vet ./...

## mod-tidy: 整理依赖
.PHONY: mod-tidy
mod-tidy:
	@echo "整理依赖..."
	$(GO) mod tidy

## mod-verify: 验证依赖
.PHONY: mod-verify
mod-verify:
	@echo "验证依赖..."
	$(GO) mod verify

## deps: 下载依赖
.PHONY: deps
deps:
	@echo "下载依赖..."
	$(GO) mod download

# ==================== Docker目标 ====================

## docker-build: 构建Docker镜像
.PHONY: docker-build
docker-build:
	@echo "构建Docker镜像..."
	docker build -t $(HarborRegistry)/$(APP_NAME):$(VERSION) -f Dockerfile .
	docker tag $(HarborRegistry)/$(APP_NAME):$(VERSION) $(HarborRegistry)/$(APP_NAME):latest
	@echo "Docker镜像构建完成: $(HarborRegistry)/$(APP_NAME):$(VERSION)"

## docker-push: 推送Docker镜像
.PHONY: docker-push
docker-push:
	@echo "推送Docker镜像..."
	docker push $(HarborRegistry)/$(APP_NAME):$(VERSION)
	docker push $(HarborRegistry)/$(APP_NAME):latest

# ==================== 部署目标 ====================

## deploy: 部署到指定环境 (env=dev|pre|pro ver=version)
.PHONY: deploy
deploy:
	@echo "部署 $(ENV).$(VER) ..."
	@cd $(DeployPath)/overlays/$(ENV) \
		&& kustomize edit add annotation ver:$(VER) -f \
		&& kustomize edit add configmap gsv-cm --behavior=merge --from-literal=ver='v-$(VER)' \
		&& kustomize edit set image $(HarborRegistry)/$(APP_NAME):$(ENV).$(VER) \
		&& cd - \
		&& kustomize build $(DeployPath)/overlays/$(ENV) | kubectl apply -f -
	@echo "部署 $(ENV).$(VER) 完成"

## build-stage: 构建并推送镜像 (env=dev|pre|pro ver=version)
.PHONY: build-stage
build-stage:
	@echo "构建镜像 $(ENV).$(VER)..."
	docker rmi -f $(HarborRegistry)/$(APP_NAME):$(ENV).$(VER) || true
	docker build --rm --no-cache -t $(HarborRegistry)/$(APP_NAME):$(ENV).$(VER) -f Dockerfile .
	docker push $(HarborRegistry)/$(APP_NAME):$(ENV).$(VER)
	@echo "构建镜像 $(ENV).$(VER) 完成"

## deploy-stage: 构建并部署 (env=dev|pre ver=BuildNum)
.PHONY: deploy-stage
deploy-stage: build-stage deploy

## help: 显示帮助信息
.PHONY: help
help:
	@echo "$(APP_NAME) Makefile 命令:"
	@echo ""
	@echo "构建和运行:"
	@echo "  make build           - 编译程序"
	@echo "  make build-linux     - 编译Linux版本"
	@echo "  make run             - 运行程序（开发环境）"
	@echo "  make run-prod        - 运行程序（生产环境）"
	@echo ""
	@echo "测试:"
	@echo "  make test            - 运行所有测试"
	@echo "  make test-coverage   - 运行测试并生成覆盖率报告"
	@echo "  make coverage        - 生成覆盖率报告（HTML）"
	@echo "  make coverage-summary - 查看覆盖率摘要"
	@echo ""
	@echo "代码质量:"
	@echo "  make fmt             - 格式化代码"
	@echo "  make vet             - 运行go vet检查"
	@echo ""
	@echo "依赖管理:"
	@echo "  make mod-tidy        - 整理依赖"
	@echo "  make mod-verify      - 验证依赖"
	@echo "  make deps            - 下载依赖"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build    - 构建Docker镜像"
	@echo "  make docker-push     - 推送Docker镜像"
	@echo ""
	@echo "部署:"
	@echo "  make build-stage     - 构建并推送镜像 (env=dev|pre|pro ver=version)"
	@echo "  make deploy-stage    - 构建并部署 (env=dev|pre ver=BuildNum)"
	@echo "  make deploy          - 部署到指定环境 (env=dev|pre|pro ver=version)"
	@echo ""
	@echo "其他:"
	@echo "  make clean           - 清理编译产物"
	@echo "  make help            - 显示此帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make run"
	@echo "  make test"
	@echo "  make deploy-stage env=dev ver=123"
