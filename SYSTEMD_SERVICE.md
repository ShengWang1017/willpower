# Willpower Forge - Systemd Service 安装指南

本指南将帮助你将 Willpower Forge 设置为 Linux 系统服务，实现开机自启动和后台运行。

## 📋 前提条件

- Linux 系统（使用 systemd，如 Ubuntu 16.04+, Debian 8+, CentOS 7+）
- 已下载 `willpower-forge-linux` 可执行文件
- sudo 权限

## 🚀 快速安装

### 步骤 1：准备应用程序目录

```bash
# 创建应用程序目录
sudo mkdir -p /opt/willpower-forge

# 复制可执行文件到目标目录
sudo cp willpower-forge-linux /opt/willpower-forge/

# 设置执行权限
sudo chmod +x /opt/willpower-forge/willpower-forge-linux

# 创建数据目录（用于存储数据库）
sudo mkdir -p /var/lib/willpower-forge

# 设置目录所有者（根据实际运行用户调整）
sudo chown -R $USER:$USER /var/lib/willpower-forge
```

### 步骤 2：修改 Service 文件

编辑 `willpower-forge.service` 文件，根据你的实际情况修改以下内容：

```ini
[Service]
# 修改为实际的用户名
User=your-username
Group=your-username

# 修改为实际的安装路径
WorkingDirectory=/var/lib/willpower-forge
ExecStart=/opt/willpower-forge/willpower-forge-linux
```

**重要配置说明：**
- `User` 和 `Group`：运行服务的用户，建议使用非 root 用户
- `WorkingDirectory`：工作目录，数据库文件将在此创建
- `ExecStart`：可执行文件的完整路径

### 步骤 3：安装 Service 文件

```bash
# 复制 service 文件到 systemd 目录
sudo cp willpower-forge.service /etc/systemd/system/

# 重新加载 systemd 配置
sudo systemctl daemon-reload

# 启用服务（开机自启动）
sudo systemctl enable willpower-forge

# 启动服务
sudo systemctl start willpower-forge
```

### 步骤 4：验证服务状态

```bash
# 查看服务状态
sudo systemctl status willpower-forge

# 查看服务日志
sudo journalctl -u willpower-forge -f

# 测试访问
curl http://localhost:8080
```

## 🔧 服务管理命令

### 基本操作

```bash
# 启动服务
sudo systemctl start willpower-forge

# 停止服务
sudo systemctl stop willpower-forge

# 重启服务
sudo systemctl restart willpower-forge

# 重新加载配置（不中断服务）
sudo systemctl reload willpower-forge

# 查看服务状态
sudo systemctl status willpower-forge

# 启用开机自启动
sudo systemctl enable willpower-forge

# 禁用开机自启动
sudo systemctl disable willpower-forge
```

### 日志查看

```bash
# 查看最新日志（实时跟踪）
sudo journalctl -u willpower-forge -f

# 查看最近 100 行日志
sudo journalctl -u willpower-forge -n 100

# 查看今天的日志
sudo journalctl -u willpower-forge --since today

# 查看指定时间范围的日志
sudo journalctl -u willpower-forge --since "2025-10-01 00:00:00" --until "2025-10-01 23:59:59"

# 查看错误日志
sudo journalctl -u willpower-forge -p err
```

## 🔐 安全配置（可选）

如果需要更严格的安全设置，可以在 service 文件的 `[Service]` 部分添加：

```ini
[Service]
# 只读系统目录
ProtectSystem=strict
ReadWritePaths=/var/lib/willpower-forge

# 禁止访问 home 目录
ProtectHome=true

# 私有 /tmp 目录
PrivateTmp=true

# 禁止提升权限
NoNewPrivileges=true

# 限制系统调用
SystemCallFilter=@system-service
SystemCallErrorNumber=EPERM

# 限制网络访问（如果只需要本地访问）
# RestrictAddressFamilies=AF_INET AF_INET6
```

修改后重新加载：
```bash
sudo systemctl daemon-reload
sudo systemctl restart willpower-forge
```

## 🌐 反向代理配置（可选）

如果需要通过域名访问，可以配置 Nginx 反向代理：

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用配置：
```bash
sudo ln -s /etc/nginx/sites-available/willpower-forge /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 🔄 更新应用程序

当有新版本时：

```bash
# 停止服务
sudo systemctl stop willpower-forge

# 备份数据库
sudo cp /var/lib/willpower-forge/willpower.db /var/lib/willpower-forge/willpower.db.backup

# 替换可执行文件
sudo cp willpower-forge-linux /opt/willpower-forge/
sudo chmod +x /opt/willpower-forge/willpower-forge-linux

# 启动服务
sudo systemctl start willpower-forge

# 查看状态确认成功
sudo systemctl status willpower-forge
```

## 🐛 故障排除

### 服务无法启动

```bash
# 查看详细错误信息
sudo systemctl status willpower-forge -l

# 查看完整日志
sudo journalctl -u willpower-forge -xe

# 检查可执行文件权限
ls -l /opt/willpower-forge/willpower-forge-linux

# 手动运行测试
cd /var/lib/willpower-forge
/opt/willpower-forge/willpower-forge-linux
```

### 端口被占用

```bash
# 查看 8080 端口占用情况
sudo lsof -i:8080
sudo netstat -tulpn | grep 8080

# 终止占用端口的进程
sudo kill -9 <PID>
```

### 权限问题

```bash
# 检查目录权限
ls -ld /var/lib/willpower-forge
ls -l /var/lib/willpower-forge/

# 修复权限
sudo chown -R your-username:your-username /var/lib/willpower-forge
sudo chmod -R 755 /var/lib/willpower-forge
```

### 数据库损坏

```bash
# 停止服务
sudo systemctl stop willpower-forge

# 恢复备份
cd /var/lib/willpower-forge
sudo cp willpower.db willpower.db.corrupted
sudo cp willpower.db.backup willpower.db

# 启动服务
sudo systemctl start willpower-forge
```

## 📁 目录结构

推荐的目录结构：

```
/opt/willpower-forge/           # 应用程序目录
├── willpower-forge-linux       # 可执行文件

/var/lib/willpower-forge/       # 数据目录
├── willpower.db                # 数据库文件
└── willpower.db.backup         # 备份文件

/etc/systemd/system/            # systemd 配置
└── willpower-forge.service     # service 文件
```

## 🔍 监控服务健康状态

创建一个简单的健康检查脚本：

```bash
#!/bin/bash
# /opt/willpower-forge/healthcheck.sh

if curl -f http://localhost:8080 >/dev/null 2>&1; then
    echo "Service is healthy"
    exit 0
else
    echo "Service is unhealthy"
    exit 1
fi
```

添加到 crontab 定期检查：
```bash
# 每 5 分钟检查一次
*/5 * * * * /opt/willpower-forge/healthcheck.sh || systemctl restart willpower-forge
```

## 📊 性能优化

如果需要处理大量请求，可以在 service 文件中添加资源限制：

```ini
[Service]
# 限制内存使用（最大 512MB）
MemoryMax=512M

# 限制 CPU 使用（最大 50%）
CPUQuota=50%

# 限制打开文件数
LimitNOFILE=4096
```

## 📝 完整示例

假设用户名为 `ubuntu`，完整的安装流程：

```bash
# 1. 创建目录
sudo mkdir -p /opt/willpower-forge
sudo mkdir -p /var/lib/willpower-forge

# 2. 复制文件
sudo cp willpower-forge-linux /opt/willpower-forge/
sudo chmod +x /opt/willpower-forge/willpower-forge-linux

# 3. 设置权限
sudo chown -R ubuntu:ubuntu /var/lib/willpower-forge

# 4. 编辑并复制 service 文件
# 修改 User=ubuntu, Group=ubuntu
sudo cp willpower-forge.service /etc/systemd/system/

# 5. 启动服务
sudo systemctl daemon-reload
sudo systemctl enable willpower-forge
sudo systemctl start willpower-forge

# 6. 验证
sudo systemctl status willpower-forge
curl http://localhost:8080
```

---

## 📮 需要帮助？

如果遇到问题：
1. 查看日志：`sudo journalctl -u willpower-forge -f`
2. 检查状态：`sudo systemctl status willpower-forge -l`
3. 提交 Issue：https://github.com/ShengWang1017/willpower/issues

---

**祝使用愉快！** 🎯
