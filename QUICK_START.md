# Willpower Forge - 快速开始指南

## 🚀 Linux 系统 - 安装为系统服务

### 前提条件
- Linux 系统（Ubuntu/Debian/CentOS 等）
- sudo 权限
- 已下载或构建好 `willpower-forge-linux` 可执行文件

### 方法一：使用自动安装脚本（推荐）

#### 步骤 1：进入项目目录
```bash
cd /home/sheng.wang/dev/willpower
```

#### 步骤 2：运行安装脚本
```bash
sudo bash install-service.sh
```

或者：
```bash
sudo ./install-service.sh
```

#### 步骤 3：按提示确认
当看到提示时，输入 `y` 并回车：
```
Configuration:
  User: your-username
  Group: your-username
  Install directory: /opt/willpower-forge
  Data directory: /var/lib/willpower-forge

Continue with installation? (y/n) y
```

#### 步骤 4：验证安装
安装完成后，会自动启动服务。你可以：

1. **查看服务状态**：
   ```bash
   sudo systemctl status willpower-forge
   ```

2. **访问应用**：
   打开浏览器访问 `http://localhost:8080`

3. **查看日志**：
   ```bash
   sudo journalctl -u willpower-forge -f
   ```

### 方法二：手动运行（不安装服务）

如果你只想临时运行，不需要安装为服务：

```bash
cd /home/sheng.wang/dev/willpower
./willpower-forge-linux
```

然后访问 `http://localhost:8080`

按 `Ctrl+C` 停止程序。

---

## 🔧 常用命令

### 服务管理
```bash
# 启动服务
sudo systemctl start willpower-forge

# 停止服务
sudo systemctl stop willpower-forge

# 重启服务
sudo systemctl restart willpower-forge

# 查看服务状态
sudo systemctl status willpower-forge

# 禁用开机自启动
sudo systemctl disable willpower-forge

# 启用开机自启动
sudo systemctl enable willpower-forge
```

### 日志查看
```bash
# 实时查看日志
sudo journalctl -u willpower-forge -f

# 查看最近 100 行日志
sudo journalctl -u willpower-forge -n 100

# 查看今天的日志
sudo journalctl -u willpower-forge --since today
```

---

## 🗑️ 卸载服务

### 卸载但保留数据
```bash
sudo bash uninstall-service.sh
```

当提示时：
- 第一个提示（确认卸载）：输入 `y`
- 第二个提示（删除数据）：输入 `n`

### 完全卸载（包括数据）
```bash
sudo bash uninstall-service.sh
```

当提示时：
- 第一个提示（确认卸载）：输入 `y`
- 第二个提示（删除数据）：输入 `y`

**注意**：选择删除数据时，脚本会先创建备份到 `/root/willpower-forge-backup-*.tar.gz`

---

## 🐛 故障排除

### 问题 1：提示 "Permission denied"
**原因**：没有执行权限

**解决**：
```bash
chmod +x install-service.sh
chmod +x willpower-forge-linux
sudo ./install-service.sh
```

### 问题 2：提示 "command not found"
**原因**：不在正确的目录

**解决**：
```bash
cd /home/sheng.wang/dev/willpower
pwd  # 确认当前目录
ls -la install-service.sh  # 确认文件存在
sudo bash install-service.sh
```

### 问题 3：服务无法启动
**查看错误信息**：
```bash
sudo journalctl -u willpower-forge -xe
```

**常见原因**：
- 端口 8080 被占用
  ```bash
  sudo lsof -i:8080
  sudo kill -9 <PID>
  ```
- 可执行文件损坏
  ```bash
  # 重新构建
  bash build-linux.sh
  ```

### 问题 4：无法访问 8080 端口
**检查防火墙**：
```bash
# Ubuntu/Debian
sudo ufw status
sudo ufw allow 8080

# CentOS/RHEL
sudo firewall-cmd --list-ports
sudo firewall-cmd --add-port=8080/tcp --permanent
sudo firewall-cmd --reload
```

---

## 📁 文件位置

### 安装后的文件位置
```
/opt/willpower-forge/              # 应用程序
├── willpower-forge-linux          # 可执行文件

/var/lib/willpower-forge/          # 数据目录
├── willpower.db                   # 数据库

/etc/systemd/system/               # 系统服务
└── willpower-forge.service        # 服务配置
```

### 源文件位置
```
/home/sheng.wang/dev/willpower/
├── install-service.sh             # 安装脚本
├── uninstall-service.sh           # 卸载脚本
├── willpower-forge-linux          # 可执行文件
├── willpower-forge.service        # 服务模板
├── build-linux.sh                 # 构建脚本
└── ...
```

---

## 🌐 远程访问

如果需要从其他机器访问，需要：

### 1. 检查防火墙（见上方"故障排除"）

### 2. 使用服务器 IP 访问
```
http://服务器IP:8080
```

### 3. 或配置 Nginx 反向代理（可选）
详见 `SYSTEMD_SERVICE.md` 中的 "反向代理配置" 部分。

---

## 📝 完整示例

假设你刚刚 clone 了项目或下载了发布版本：

```bash
# 1. 进入项目目录
cd /home/sheng.wang/dev/willpower

# 2. 确认文件存在
ls -la willpower-forge-linux install-service.sh

# 3. 如果需要，添加执行权限
chmod +x install-service.sh willpower-forge-linux

# 4. 运行安装脚本
sudo ./install-service.sh

# 5. 输入 y 确认安装

# 6. 查看服务状态
sudo systemctl status willpower-forge

# 7. 访问应用
# 打开浏览器访问 http://localhost:8080

# 8. 查看实时日志
sudo journalctl -u willpower-forge -f
```

---

## ❓ 还有问题？

1. 查看详细文档：`SYSTEMD_SERVICE.md`
2. 查看 README：`README.md`
3. 提交 Issue：https://github.com/ShengWang1017/willpower/issues

---

**祝使用愉快！** 🎯
