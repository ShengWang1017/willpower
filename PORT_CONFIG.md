# 端口配置指南

Willpower Forge 支持自定义监听端口，默认端口为 **5173**。

## 📋 修改端口的方法

### 方法一：修改 systemd service 环境变量（推荐）

如果使用 systemd 服务运行：

#### 1. 编辑 service 文件

```bash
sudo systemctl edit --full willpower-forge
```

或直接编辑：
```bash
sudo nano /etc/systemd/system/willpower-forge.service
```

#### 2. 修改 Environment 变量

找到 `Environment="PORT=5173"` 这一行，修改为你想要的端口：

```ini
[Service]
# ... 其他配置 ...

# Environment variables
# Set PORT to change the listening port (default: 5173)
Environment="PORT=8080"     # 改成你想要的端口，如 8080

# ... 其他配置 ...
```

#### 3. 重新加载并重启服务

```bash
sudo systemctl daemon-reload
sudo systemctl restart willpower-forge
```

#### 4. 验证端口

```bash
# 查看服务状态
sudo systemctl status willpower-forge

# 查看监听端口
sudo ss -tulpn | grep willpower-forge

# 查看日志确认端口
sudo journalctl -u willpower-forge -n 20
```

应该能看到：
```
Starting server on 0.0.0.0:8080
```

---

### 方法二：设置环境变量（临时运行）

如果直接运行可执行文件：

```bash
# 使用默认端口 5173
./willpower-forge-linux

# 自定义端口
PORT=8080 ./willpower-forge-linux

# 或者
export PORT=8080
./willpower-forge-linux
```

---

### 方法三：修改源代码（永久修改默认值）

如果你想修改默认端口：

#### 1. 编辑 main.go

```bash
nano willpower-forge-api/main.go
```

#### 2. 修改默认端口

找到这几行代码：
```go
// Get port from environment variable, default to 5173
port := os.Getenv("PORT")
if port == "" {
    port = "5173"  // 修改这里的默认端口
}
```

改为：
```go
// Get port from environment variable, default to 8080
port := os.Getenv("PORT")
if port == "" {
    port = "8080"  // 你的新默认端口
}
```

#### 3. 重新构建

**Linux:**
```bash
cd willpower-forge-api
go build -ldflags="-s -w" -o ../willpower-forge-linux .
cd ..
```

**Windows (在 Windows 上运行):**
```cmd
cd willpower-forge-api
go build -ldflags="-s -w" -o ..\willpower-forge-windows.exe .
cd ..
```

#### 4. 重新安装服务

```bash
sudo ./uninstall-service.sh
sudo ./install-service.sh
```

---

## 🔄 快速切换端口示例

### 从 5173 改为 8080

```bash
# 1. 编辑服务文件
sudo nano /etc/systemd/system/willpower-forge.service

# 2. 修改这一行：
Environment="PORT=8080"

# 3. 重新加载并重启
sudo systemctl daemon-reload
sudo systemctl restart willpower-forge

# 4. 验证
sudo journalctl -u willpower-forge -n 5
```

### 从 8080 改为 3000

```bash
sudo systemctl stop willpower-forge
sudo sed -i 's/PORT=8080/PORT=3000/g' /etc/systemd/system/willpower-forge.service
sudo systemctl daemon-reload
sudo systemctl start willpower-forge
sudo systemctl status willpower-forge
```

---

## 🔍 验证端口配置

### 1. 检查服务日志

```bash
sudo journalctl -u willpower-forge -f
```

应该看到：
```
Starting server on 0.0.0.0:YOUR_PORT
[GIN-debug] Listening and serving HTTP on 0.0.0.0:YOUR_PORT
```

### 2. 检查监听端口

```bash
sudo ss -tulpn | grep YOUR_PORT
```

应该看到：
```
tcp   LISTEN 0      4096    *:YOUR_PORT    *:*    users:(("willpower-forge",pid=XXXX,fd=7))
```

### 3. 测试访问

```bash
# 本地测试
curl http://localhost:YOUR_PORT

# 查看公网 IP
curl ifconfig.me

# 外网访问（需要配置安全组）
# http://你的公网IP:YOUR_PORT
```

---

## 🔒 安全组配置

**重要：修改端口后，需要更新云服务器安全组规则！**

### 阿里云 ECS

1. 登录阿里云控制台
2. 进入 ECS 实例 → 安全组
3. **删除**旧端口规则（如 8080）
4. **添加**新端口规则（如 5173）：
   ```
   协议类型：TCP
   端口范围：5173/5173
   授权对象：0.0.0.0/0
   描述：Willpower Forge
   ```

### 其他云服务商

参考 `EXTERNAL_ACCESS.md` 文档。

---

## 🎯 常见端口选择

- **5173**: Vite 默认开发端口
- **8080**: 常用的 Web 应用端口
- **3000**: Node.js 应用常用端口
- **80**: HTTP 标准端口（需要 root 权限）
- **443**: HTTPS 标准端口（需要 root 权限和 SSL 证书）

### 使用 80 或 443 端口

这些是特权端口，需要特殊配置：

#### 选项 1：使用 authbind

```bash
# 安装 authbind
sudo apt-get install authbind

# 允许使用 80 端口
sudo touch /etc/authbind/byport/80
sudo chmod 500 /etc/authbind/byport/80
sudo chown your-username /etc/authbind/byport/80

# 修改 service 文件
sudo nano /etc/systemd/system/willpower-forge.service
```

修改 `ExecStart` 行：
```ini
ExecStart=/usr/bin/authbind --deep /opt/willpower-forge/willpower-forge-linux
Environment="PORT=80"
```

#### 选项 2：使用 setcap（不推荐）

```bash
sudo setcap CAP_NET_BIND_SERVICE=+eip /opt/willpower-forge/willpower-forge-linux
```

#### 选项 3：使用 Nginx 反向代理（推荐）

让 Nginx 监听 80/443，然后转发到 5173：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:5173;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

这样你可以：
- Willpower Forge 使用 5173 端口（无需特权）
- Nginx 监听 80/443 端口（提供 HTTPS、负载均衡等功能）

---

## 🐛 故障排除

### 端口被占用

```bash
# 查看占用端口的进程
sudo lsof -i:5173

# 或
sudo ss -tulpn | grep 5173

# 结束占用进程
sudo kill -9 <PID>
```

### 端口修改不生效

```bash
# 确认 service 文件已修改
sudo cat /etc/systemd/system/willpower-forge.service | grep PORT

# 确认已重新加载
sudo systemctl daemon-reload

# 确认已重启
sudo systemctl restart willpower-forge

# 查看实际监听端口
sudo ss -tulpn | grep willpower-forge
```

### 防火墙阻止

```bash
# Ubuntu/Debian - ufw
sudo ufw allow 5173/tcp
sudo ufw status

# CentOS/RHEL - firewalld
sudo firewall-cmd --add-port=5173/tcp --permanent
sudo firewall-cmd --reload
```

---

## 📝 配置示例

### 示例 1：开发环境（5173）

```ini
# /etc/systemd/system/willpower-forge.service
[Service]
Environment="PORT=5173"
Environment="GIN_MODE=debug"
```

### 示例 2：生产环境（8080）

```ini
# /etc/systemd/system/willpower-forge.service
[Service]
Environment="PORT=8080"
Environment="GIN_MODE=release"
```

### 示例 3：多实例运行

如果需要运行多个实例：

```bash
# 实例 1 - 端口 5173
sudo cp /etc/systemd/system/willpower-forge.service /etc/systemd/system/willpower-forge-1.service
sudo nano /etc/systemd/system/willpower-forge-1.service
# 修改 PORT=5173, WorkingDirectory=/var/lib/willpower-forge-1

# 实例 2 - 端口 5174
sudo cp /etc/systemd/system/willpower-forge.service /etc/systemd/system/willpower-forge-2.service
sudo nano /etc/systemd/system/willpower-forge-2.service
# 修改 PORT=5174, WorkingDirectory=/var/lib/willpower-forge-2

# 启动
sudo systemctl daemon-reload
sudo systemctl start willpower-forge-1
sudo systemctl start willpower-forge-2
```

---

## ✅ 配置检查清单

修改端口后，请检查：

- [ ] service 文件中的 `Environment="PORT=XXXX"` 已修改
- [ ] 已执行 `systemctl daemon-reload`
- [ ] 已执行 `systemctl restart willpower-forge`
- [ ] 日志中显示正确的端口：`journalctl -u willpower-forge -n 5`
- [ ] 端口正在监听：`ss -tulpn | grep XXXX`
- [ ] 本地可以访问：`curl http://localhost:XXXX`
- [ ] 云服务器安全组已更新
- [ ] 外网可以访问：`http://公网IP:XXXX`

---

## 📞 需要帮助？

如果仍有问题：
1. 查看服务状态：`sudo systemctl status willpower-forge -l`
2. 查看详细日志：`sudo journalctl -u willpower-forge -xe`
3. 提交 Issue：https://github.com/ShengWang1017/willpower/issues

---

**祝配置顺利！** 🚀
