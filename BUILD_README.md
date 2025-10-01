# Willpower Forge - 打包说明

本项目已配置为可打包成单一可执行文件，包含前后端所有功能。

## 📦 已生成的可执行文件

当前项目根目录包含以下可执行文件：

- **willpower-forge-windows.exe** (15MB) - Windows 版本
- **willpower-forge-linux** (18MB) - Linux 版本

这些文件已经嵌入了完整的前端界面，可以直接运行。

## 🚀 快速使用

### Windows 系统
1. 双击 `willpower-forge-windows.exe`
2. 打开浏览器访问 `http://localhost:5173`

### Linux 系统
```bash
chmod +x willpower-forge-linux
./willpower-forge-linux
```
然后在浏览器访问 `http://localhost:5173`

## 🔨 重新打包

如果你修改了代码，需要重新打包，可以使用以下脚本：

### Windows 系统
双击运行 `build-windows.bat` 或在命令行执行：
```cmd
build-windows.bat
```

### Linux/Mac 系统
```bash
chmod +x build-linux.sh
./build-linux.sh
```

## 📋 打包脚本功能

两个打包脚本会自动执行以下步骤：

1. 安装前端依赖（npm install）
2. 构建前端静态文件（npm run build）
3. 复制前端文件到后端目录
4. 编译 Go 程序并嵌入前端文件
5. 生成单一可执行文件

## 🔧 环境要求

只有在**重新打包**时才需要以下环境：

- Node.js (v14+) 和 npm
- Go (v1.18+)

**运行**已打包的可执行文件不需要任何环境依赖！

## 📝 注意事项

1. **数据库文件**：首次运行会在可执行文件同目录生成 `willpower.db` 数据库文件
2. **端口配置**：默认使用 5173 端口，可通过环境变量 `PORT` 自定义端口（如 `PORT=8080 ./willpower-forge-linux`）
3. **跨平台**：Windows 可执行文件只能在 Windows 上运行，Linux 可执行文件只能在 Linux 上运行
4. **文件大小**：可执行文件较大是因为包含了完整的前后端代码和依赖

## 🎯 技术实现

使用 Go 1.16+ 的 `embed` 功能将前端静态文件直接嵌入到二进制文件中：

```go
//go:embed web/dist
var webFS embed.FS
```

这样就实现了真正的"单文件部署"，无需额外的文件或依赖。

## 📁 项目结构

```
willpower/
├── willpower-forge-api/        # Go 后端
│   ├── web/dist/               # 嵌入的前端文件
│   └── ...
├── willpower-forge-web/        # Vue 前端
│   └── dist/                   # 构建输出
├── build-windows.bat           # Windows 打包脚本
├── build-linux.sh              # Linux 打包脚本
├── willpower-forge-windows.exe # Windows 可执行文件
└── willpower-forge-linux       # Linux 可执行文件
```

## 🐛 故障排除

### 端口被占用
如果提示 5173 端口被占用，可以：
- 关闭占用该端口的程序，或
- 使用自定义端口：`PORT=8080 ./willpower-forge-linux`

### Windows 防火墙警告
首次运行可能会弹出防火墙警告，点击"允许访问"即可。

### Linux 权限问题
如果提示权限不足，执行：
```bash
chmod +x willpower-forge-linux
```

## 📮 分发说明

要将应用分发给其他用户，只需：

1. 复制对应平台的可执行文件
2. 告知用户直接运行即可
3. 无需安装任何运行时环境

用户计算机上不需要安装 Node.js、Go 或任何其他依赖！
