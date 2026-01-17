#!/bin/bash
# 部署脚本
#
# 用于构建和部署 Game Slots Version Server

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    log_info "Docker 已安装: $(docker --version)"
}

# 检查 kubectl 是否安装
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        log_warn "kubectl 未安装，跳过 Kubernetes 部署"
        return 1
    fi
    log_info "kubectl 已安装: $(kubectl version --client --short 2>&1)"
    return 0
}

# 构建 Docker 镜像
build_image() {
    log_info "构建 Docker 镜像..."
    docker build -f deployments/docker/Dockerfile -t game-slots-vsn:latest .
    log_info "Docker 镜像构建完成"
}

# 运行本地开发环境
run_dev() {
    log_info "启动本地开发环境..."
    cd deployments/docker
    docker-compose up -d
    log_info "服务已启动"
    log_info "访问地址: http://localhost:8080"
    log_info "查看日志: docker-compose logs -f"
}

# 停止本地开发环境
stop_dev() {
    log_info "停止本地开发环境..."
    cd deployments/docker
    docker-compose down
    log_info "服务已停止"
}

# 部署到 Kubernetes
deploy_k8s() {
    if ! check_kubectl; then
        return
    fi

    log_info "部署到 Kubernetes..."

    # 创建 namespace（如果不存在）
    kubectl create namespace game-slots-vsn --dry-run=client -o yaml | kubectl apply -f -

    # 应用 ConfigMap
    log_info "应用 ConfigMap..."
    kubectl apply -f deployments/k8s/configmap.yaml

    # 应用 Secret（需要先修改密码）
    log_warn "请先修改 deployments/k8s/secret.yaml.template 中的密码"
    if [ -f "deployments/k8s/secret.yaml" ]; then
        kubectl apply -f deployments/k8s/secret.yaml
    else
        log_error "Secret 文件不存在，请从 template 创建并修改密码"
        return
    fi

    # 应用 Deployment
    log_info "应用 Deployment..."
    kubectl apply -f deployments/k8s/deployment.yaml

    # 等待 Pod 就绪
    log_info "等待 Pod 就绪..."
    kubectl wait --for=condition=available --timeout=60s \
        deployment/game-slots-vsn -n default

    log_info "部署完成"
    log_info "查看状态: kubectl get pods -l app=game-slots-vsn"
    log_info "查看日志: kubectl logs -l app=game-slots-vsn -f"
}

# 从 Kubernetes 删除
delete_k8s() {
    if ! check_kubectl; then
        return
    fi

    log_info "从 Kubernetes 删除..."
    kubectl delete -f deployments/k8s/deployment.yaml --ignore-not-found=true
    kubectl delete -f deployments/k8s/configmap.yaml --ignore-not-found=true
    kubectl delete -f deployments/k8s/secret.yaml --ignore-not-found=true
    log_info "删除完成"
}

# 运行测试
run_tests() {
    log_info "运行测试..."
    go test -v -race -coverprofile=coverage.out ./...
    log_info "测试完成"
}

# 运行 benchmark
run_benchmark() {
    log_info "运行 benchmark..."
    go test -bench=. -benchmem ./...
    log_info "benchmark 完成"
}

# 显示帮助信息
show_help() {
    echo "Game Slots Version Server 部署脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 <command>"
    echo ""
    echo "可用命令:"
    echo "  build        构建 Docker 镜像"
    echo "  dev          启动本地开发环境"
    echo "  stop-dev     停止本地开发环境"
    echo "  deploy-k8s   部署到 Kubernetes"
    echo "  delete-k8s   从 Kubernetes 删除"
    echo "  test         运行测试"
    echo "  benchmark    运行 benchmark"
    echo "  help         显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 build        # 构建镜像"
    echo "  $0 dev          # 启动开发环境"
    echo "  $0 deploy-k8s   # 部署到 K8s"
}

# 主函数
main() {
    case "$1" in
        build)
            check_docker
            build_image
            ;;
        dev)
            check_docker
            run_dev
            ;;
        stop-dev)
            check_docker
            stop_dev
            ;;
        deploy-k8s)
            check_docker
            deploy_k8s
            ;;
        delete-k8s)
            check_kubectl && delete_k8s
            ;;
        test)
            run_tests
            ;;
        benchmark)
            run_benchmark
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
