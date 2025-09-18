# ToDoList 项目 Makefile
# 集成前后端开发、构建、测试和部署的所有命令

.PHONY: help install dev build test clean proto start stop health check-deps

# 默认目标
.DEFAULT_GOAL := help

# 项目配置
BACKEND_DIR := backend
FRONTEND_DIR := frontend
PROTO_DIR := $(BACKEND_DIR)/proto
GO_MODULE := github.com/jjeejj/todolist/backend

# 端口配置
BACKEND_PORT := 8080
FRONTEND_PORT := 3000

# 颜色输出
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

## 帮助信息
help: ## 显示帮助信息
	@echo "$(BLUE)ToDoList 项目 Makefile 命令$(RESET)"
	@echo ""
	@echo "$(GREEN)开发命令:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-15s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## 依赖检查和安装
check-deps: ## 检查必要的依赖工具
	@echo "$(BLUE)检查依赖工具...$(RESET)"
	@command -v go >/dev/null 2>&1 || { echo "$(RED)错误: 需要安装 Go$(RESET)"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "$(RED)错误: 需要安装 Node.js$(RESET)"; exit 1; }
	@command -v npm >/dev/null 2>&1 || { echo "$(RED)错误: 需要安装 npm$(RESET)"; exit 1; }
	@command -v protoc >/dev/null 2>&1 || { echo "$(RED)错误: 需要安装 protoc$(RESET)"; exit 1; }
	@echo "$(GREEN)✓ 所有依赖工具已安装$(RESET)"

install: check-deps ## 安装所有依赖
	@echo "$(BLUE)安装后端依赖...$(RESET)"
	cd $(BACKEND_DIR) && go mod download
	cd $(BACKEND_DIR) && go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	@echo "$(BLUE)安装前端依赖...$(RESET)"
	cd $(FRONTEND_DIR) && npm install
	@echo "$(GREEN)✓ 依赖安装完成$(RESET)"

## Protobuf 相关
proto-gen: ## 生成 protobuf 代码
	@echo "$(BLUE)生成后端 protobuf 代码...$(RESET)"
	cd $(BACKEND_DIR) && protoc \
		--go_out=. --go_opt=paths=source_relative \
		--connect-go_out=. --connect-go_opt=paths=source_relative \
		proto/todolist/v1/todolist.proto
	@echo "$(BLUE)生成前端 protobuf 代码...$(RESET)"
	cd $(FRONTEND_DIR) && npx buf generate ../$(BACKEND_DIR)/proto
	@echo "$(GREEN)✓ Protobuf 代码生成完成$(RESET)"

proto-clean: ## 清理生成的 protobuf 代码
	@echo "$(BLUE)清理 protobuf 生成文件...$(RESET)"
	rm -f $(BACKEND_DIR)/proto/todolist/v1/*.pb.go
	rm -f $(BACKEND_DIR)/proto/todolist/v1/*connect/*.go
	rm -rf $(FRONTEND_DIR)/src/gen
	@echo "$(GREEN)✓ Protobuf 文件清理完成$(RESET)"

## 开发服务
dev-backend: ## 启动后端开发服务器
	@echo "$(BLUE)启动后端服务器 (端口 $(BACKEND_PORT))...$(RESET)"
	cd $(BACKEND_DIR) && go run cmd/server/main.go

dev-frontend: ## 启动前端开发服务器
	@echo "$(BLUE)启动前端服务器 (端口 $(FRONTEND_PORT))...$(RESET)"
	cd $(FRONTEND_DIR) && npm run dev

dev: ## 同时启动前后端开发服务器
	@echo "$(BLUE)启动完整开发环境...$(RESET)"
	@echo "$(YELLOW)注意: 请在不同终端窗口中分别运行:$(RESET)"
	@echo "  make dev-backend  # 启动后端服务"
	@echo "  make dev-frontend # 启动前端服务"

## 构建
build-backend: ## 构建后端应用
	@echo "$(BLUE)构建后端应用...$(RESET)"
	cd $(BACKEND_DIR) && go build -o bin/server cmd/server/main.go
	@echo "$(GREEN)✓ 后端构建完成: $(BACKEND_DIR)/bin/server$(RESET)"

build-frontend: ## 构建前端应用
	@echo "$(BLUE)构建前端应用...$(RESET)"
	cd $(FRONTEND_DIR) && npm run build
	@echo "$(GREEN)✓ 前端构建完成$(RESET)"

build: proto-gen build-backend build-frontend ## 构建完整应用

## 测试
test-backend: ## 运行后端测试
	@echo "$(BLUE)运行后端测试...$(RESET)"
	cd $(BACKEND_DIR) && go test -v ./...

test-frontend: ## 运行前端测试
	@echo "$(BLUE)运行前端测试...$(RESET)"
	cd $(FRONTEND_DIR) && npm test

test: test-backend test-frontend ## 运行所有测试

## 代码质量
lint-backend: ## 后端代码检查
	@echo "$(BLUE)检查后端代码...$(RESET)"
	cd $(BACKEND_DIR) && go vet ./...
	cd $(BACKEND_DIR) && go fmt ./...

lint-frontend: ## 前端代码检查
	@echo "$(BLUE)检查前端代码...$(RESET)"
	cd $(FRONTEND_DIR) && npm run lint

lint: lint-backend lint-frontend ## 运行所有代码检查

## 生产部署
start-backend: build-backend ## 启动后端生产服务
	@echo "$(BLUE)启动后端生产服务...$(RESET)"
	cd $(BACKEND_DIR) && ./bin/server

start-frontend: build-frontend ## 启动前端生产服务
	@echo "$(BLUE)启动前端生产服务...$(RESET)"
	cd $(FRONTEND_DIR) && npm start

start: start-backend start-frontend ## 启动生产服务

## 健康检查
health: ## 检查服务健康状态
	@echo "$(BLUE)检查后端服务健康状态...$(RESET)"
	@curl -f http://localhost:$(BACKEND_PORT)/health >/dev/null 2>&1 && \
		echo "$(GREEN)✓ 后端服务正常$(RESET)" || \
		echo "$(RED)✗ 后端服务异常$(RESET)"
	@echo "$(BLUE)检查前端服务健康状态...$(RESET)"
	@curl -f http://localhost:$(FRONTEND_PORT) >/dev/null 2>&1 && \
		echo "$(GREEN)✓ 前端服务正常$(RESET)" || \
		echo "$(RED)✗ 前端服务异常$(RESET)"

## 清理
clean-backend: ## 清理后端构建文件
	@echo "$(BLUE)清理后端构建文件...$(RESET)"
	rm -rf $(BACKEND_DIR)/bin
	cd $(BACKEND_DIR) && go clean

clean-frontend: ## 清理前端构建文件
	@echo "$(BLUE)清理前端构建文件...$(RESET)"
	rm -rf $(FRONTEND_DIR)/.next
	rm -rf $(FRONTEND_DIR)/out
	rm -rf $(FRONTEND_DIR)/node_modules/.cache

clean: clean-backend clean-frontend proto-clean ## 清理所有构建文件

## 重置项目
reset: clean ## 完全重置项目
	@echo "$(BLUE)重置项目...$(RESET)"
	rm -rf $(FRONTEND_DIR)/node_modules
	cd $(BACKEND_DIR) && go clean -modcache
	@echo "$(GREEN)✓ 项目重置完成$(RESET)"

## 快速启动
quick-start: install proto-gen ## 快速启动开发环境
	@echo "$(GREEN)✓ 开发环境准备完成!$(RESET)"
	@echo ""
	@echo "$(YELLOW)现在可以运行以下命令启动服务:$(RESET)"
	@echo "  make dev-backend  # 启动后端服务 (端口 $(BACKEND_PORT))"
	@echo "  make dev-frontend # 启动前端服务 (端口 $(FRONTEND_PORT))"
	@echo ""
	@echo "$(YELLOW)或者查看所有可用命令:$(RESET)"
	@echo "  make help"

## 项目信息
info: ## 显示项目信息
	@echo "$(BLUE)项目信息:$(RESET)"
	@echo "  项目名称: ToDoList"
	@echo "  后端目录: $(BACKEND_DIR)"
	@echo "  前端目录: $(FRONTEND_DIR)"
	@echo "  后端端口: $(BACKEND_PORT)"
	@echo "  前端端口: $(FRONTEND_PORT)"
	@echo "  Go 模块: $(GO_MODULE)"
	@echo ""
	@echo "$(BLUE)服务地址:$(RESET)"
	@echo "  后端 API: http://localhost:$(BACKEND_PORT)"
	@echo "  前端应用: http://localhost:$(FRONTEND_PORT)"
	@echo "  健康检查: http://localhost:$(BACKEND_PORT)/health"