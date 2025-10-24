# InvesTool 前端调试指南

## 问题解决

刚才遇到的错误是因为缺少 `public/index.html` 文件，现在已经创建了必要的文件。

## 启动步骤

### 1. 确保后端服务运行
```bash
# 在项目根目录启动后端
go run main.go webserver
```

### 2. 启动前端开发服务器
```bash
# 进入前端目录
cd frontend

# 安装依赖（如果还没安装）
npm install

# 启动开发服务器
npm run dev
```

### 3. 访问应用
打开浏览器访问：http://localhost:3000

## 可用的调试命令

```bash
npm run dev          # 启动开发服务器
npm run dev:debug    # 启用源码映射的调试模式
npm run dev:port     # 使用 3001 端口启动
npm run build        # 构建生产版本
npm run test         # 运行测试
npm run lint         # 代码检查
npm run lint:fix     # 自动修复代码格式
npm run type-check   # TypeScript 类型检查
```

## 调试功能

- ✅ 热重载：修改代码自动刷新
- ✅ 错误提示：实时显示编译错误
- ✅ 源码映射：支持断点调试
- ✅ 代理配置：自动代理 API 请求到后端 (http://localhost:8080)

## 常见问题

### 1. 端口冲突
如果 3000 端口被占用，使用：
```bash
npm run dev:port
```

### 2. API 请求失败
- 确保后端服务在 8080 端口运行
- 检查 CORS 配置
- 查看浏览器 Network 面板

### 3. 依赖问题
```bash
# 清除缓存重新安装
rm -rf node_modules package-lock.json
npm install
```

现在应该可以正常启动调试了！
