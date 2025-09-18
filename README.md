# 📝 ToDoList - 待办事项管理应用

一个基于 Go + Next.js + ConnectRPC 构建的现代化待办事项管理应用，提供简洁高效的任务管理体验。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)
![Node Version](https://img.shields.io/badge/node-18+-green.svg)

## ✨ 特性

- 🚀 **现代化技术栈**: Go + Next.js + ConnectRPC
- 📱 **响应式设计**: 支持桌面和移动端
- ⚡ **实时更新**: 基于 HTTP/2 的高性能通信
- 🎨 **美观界面**: TailwindCSS + Lucide Icons
- 🔧 **类型安全**: TypeScript + Protocol Buffers
- 🛠️ **开发友好**: 完整的 Makefile 工具链

## 🏗️ 技术架构

```
┌─────────────────┐    ConnectRPC     ┌─────────────────┐
│   Next.js 前端   │ ◄──────────────► │   Go 后端服务    │
│                 │    (HTTP/2)       │                 │
│ • React 19      │                   │ • ConnectRPC    │
│ • TypeScript    │                   │ • Protocol Buf  │
│ • TailwindCSS   │                   │ • 内存存储       │
└─────────────────┘                   └─────────────────┘
```

### 技术栈

**前端**
- [Next.js 15.5.3](https://nextjs.org/) - React 全栈框架
- [React 19](https://react.dev/) - 用户界面库
- [TypeScript 5](https://www.typescriptlang.org/) - 类型安全
- [TailwindCSS 4](https://tailwindcss.com/) - 原子化 CSS 框架
- [Lucide React](https://lucide.dev/) - 图标库
- [ConnectRPC](https://connectrpc.com/) - 类型安全的 RPC 客户端

**后端**
- [Go 1.21+](https://golang.org/) - 高性能后端语言
- [ConnectRPC](https://connectrpc.com/) - 基于 HTTP/2 的 RPC 框架
- [Protocol Buffers](https://protobuf.dev/) - 接口定义和序列化
- 内存存储 - 轻量级数据存储

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- npm 或 yarn
- protoc (Protocol Buffers 编译器)

### 安装依赖

```bash
# 检查环境依赖
make check-deps

# 安装所有依赖
make install

# 生成 protobuf 代码
make proto-gen
```

### 启动开发服务

```bash
# 方式一：使用 Makefile（推荐）
make dev-backend   # 启动后端服务 (端口 8080)
make dev-frontend  # 启动前端服务 (端口 3000)

# 方式二：手动启动
cd backend && go run cmd/server/main.go
cd frontend && npm run dev
```

### 访问应用

- 🌐 **前端应用**: http://localhost:3000
- 🔧 **后端 API**: http://localhost:8080
- 💚 **健康检查**: http://localhost:8080/health

## 📖 使用说明

### 基本功能

1. **添加任务**: 在输入框中输入任务内容，点击"添加任务"按钮
2. **查看任务**: 所有任务会实时显示在任务列表中
3. **删除任务**: 点击任务右侧的删除按钮移除任务

### 界面预览

```
┌─────────────────────────────────────────┐
│              待办事项管理                │
├─────────────────────────────────────────┤
│ [输入新任务...        ] [添加任务]      │
├─────────────────────────────────────────┤
│ □ 完成项目文档                    [删除] │
│ □ 准备会议材料                    [删除] │
│ □ 回复客户邮件                    [删除] │
└─────────────────────────────────────────┘
```

## 🛠️ 开发指南

### 项目结构

```
ToDoList/
├── backend/                 # Go 后端服务
│   ├── cmd/server/         # 服务器入口
│   ├── internal/           # 内部包
│   │   ├── repository/     # 数据访问层
│   │   └── service/        # 业务逻辑层
│   └── proto/              # Protocol Buffers 定义
├── frontend/               # Next.js 前端应用
│   ├── src/
│   │   ├── app/           # Next.js App Router
│   │   ├── components/    # React 组件
│   │   ├── gen/           # 生成的 protobuf 代码
│   │   └── lib/           # 工具库
│   └── public/            # 静态资源
└── Makefile               # 构建工具
```

### 常用命令

```bash
# 开发相关
make help          # 查看所有可用命令
make quick-start   # 快速启动开发环境
make dev-backend   # 启动后端开发服务
make dev-frontend  # 启动前端开发服务

# 构建和测试
make build         # 构建完整应用
make test          # 运行所有测试
make lint          # 代码检查

# Protobuf 管理
make proto-gen     # 生成 protobuf 代码
make proto-clean   # 清理生成的代码

# 项目维护
make clean         # 清理构建文件
make reset         # 完全重置项目
make health        # 检查服务状态
```

### API 接口

应用使用 ConnectRPC 提供以下 API：

#### 添加任务
```protobuf
rpc AddTask(AddTaskRequest) returns (AddTaskResponse)
```

#### 获取任务列表
```protobuf
rpc GetTasks(GetTasksRequest) returns (GetTasksResponse)
```

#### 删除任务
```protobuf
rpc DeleteTask(DeleteTaskRequest) returns (DeleteTaskResponse)
```

### 数据模型

```protobuf
message Task {
  string id = 1;           // 任务唯一标识
  string text = 2;         // 任务内容
  int64 created_at = 3;    // 创建时间戳
}
```

## 🧪 测试

```bash
# 运行后端测试
make test-backend

# 运行前端测试
make test-frontend

# 运行所有测试
make test
```

## 📦 构建部署

### 开发环境构建

```bash
make build
```

### 生产环境部署

```bash
# 构建生产版本
make build

# 启动生产服务
make start-backend   # 后端服务
make start-frontend  # 前端服务
```

## 🔧 配置说明

### 端口配置

- 后端服务: `8080`
- 前端服务: `3000`

可在 `Makefile` 中修改端口配置：

```makefile
BACKEND_PORT := 8080
FRONTEND_PORT := 3000
```

### 环境变量

当前版本使用默认配置，无需额外环境变量。

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 开发规范

- 遵循 Go 和 TypeScript 的最佳实践
- 添加适当的测试覆盖
- 更新相关文档
- 确保代码通过 lint 检查

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [ConnectRPC](https://connectrpc.com/) - 优秀的 RPC 框架
- [Next.js](https://nextjs.org/) - 强大的 React 框架
- [TailwindCSS](https://tailwindcss.com/) - 实用的 CSS 框架
- [Go](https://golang.org/) - 高效的后端语言

## 📞 支持

如果你有任何问题或建议，请：

- 提交 [Issue](../../issues)
- 发起 [Discussion](../../discussions)
- 查看 [Wiki](../../wiki) 文档

---

⭐ 如果这个项目对你有帮助，请给它一个星标！