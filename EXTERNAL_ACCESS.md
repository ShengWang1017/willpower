# 允许外网访问 Willpower Forge

## 问题诊断

如果你的服务已经启动，但外网无法访问，通常是以下原因之一：

### ✅ 1. 服务监听地址（已正确配置）
服务已配置为监听 `0.0.0.0:8080`，可以接受来自任何网络接口的连接。

### ✅ 2. 本地防火墙（Ubuntu 默认未启用）
Ubuntu 默认的 ufw 防火墙通常是关闭的，不会阻止连接。

### ❌ 3. 云服务器安全组（需要配置）
**这是最常见的问题！** 云服务器默认只开放 22(SSH) 和 3389(RDP) 端口。

---

## 🔓 解决方案：配置云服务器安全组

根据你使用的云服务商，按以下步骤操作：

## 阿里云 ECS

### 方法一：通过控制台配置

1. **登录阿里云控制台**
   - 访问：https://ecs.console.aliyun.com/

2. **进入实例管理**
   - 点击左侧菜单：**实例与镜像** > **实例**
   - 找到你的 ECS 实例

3. **配置安全组**
   - 点击实例 ID 进入详情页
   - 点击 **安全组** 选项卡
   - 点击安全组 ID 进入安全组规则页面

4. **添加入方向规则**
   - 点击 **配置规则**
   - 点击 **入方向** > **手动添加**
   - 填写规则：
     ```
     授权策略：允许
     优先级：1
     协议类型：自定义 TCP
     端口范围：8080/8080
     授权对象：0.0.0.0/0
     描述：Willpower Forge 应用端口
     ```
   - 点击 **保存**

5. **验证访问**
   - 打开浏览器访问：`http://你的公网IP:8080`

### 方法二：使用阿里云 CLI

```bash
# 安装阿里云 CLI
wget https://aliyuncli.alicdn.com/aliyun-cli-linux-latest-amd64.tgz
tar -xzf aliyun-cli-linux-latest-amd64.tgz
sudo mv aliyun /usr/local/bin/

# 配置凭证
aliyun configure

# 添加安全组规则
aliyun ecs AuthorizeSecurityGroup \
  --SecurityGroupId sg-xxxxxx \
  --IpProtocol tcp \
  --PortRange 8080/8080 \
  --SourceCidrIp 0.0.0.0/0 \
  --Description "Willpower Forge"
```

### 快速查找你的公网 IP

```bash
# 在服务器上执行
curl ifconfig.me
# 或
curl ipinfo.io/ip
```

---

## 腾讯云 CVM

1. **登录腾讯云控制台**
   - 访问：https://console.cloud.tencent.com/cvm

2. **进入实例管理**
   - 点击左侧 **云服务器**
   - 找到你的实例

3. **配置安全组**
   - 点击实例 ID
   - 点击 **安全组** 选项卡
   - 点击 **编辑规则**

4. **添加入站规则**
   ```
   类型：自定义
   来源：0.0.0.0/0
   协议端口：TCP:8080
   策略：允许
   备注：Willpower Forge
   ```

5. **保存并验证**

---

## AWS EC2

1. **登录 AWS 控制台**
   - 访问：https://console.aws.amazon.com/ec2/

2. **进入实例详情**
   - 选择你的 EC2 实例
   - 点击 **Security** 选项卡
   - 点击 Security Group 链接

3. **编辑入站规则**
   - 点击 **Edit inbound rules**
   - 点击 **Add rule**
   ```
   Type: Custom TCP
   Port range: 8080
   Source: 0.0.0.0/0
   Description: Willpower Forge
   ```

4. **保存规则**

---

## Google Cloud (GCP)

1. **登录 GCP 控制台**
   - 访问：https://console.cloud.google.com/

2. **配置防火墙规则**
   - 导航到：**VPC 网络** > **防火墙**
   - 点击 **创建防火墙规则**

3. **创建规则**
   ```
   名称：allow-willpower-forge
   目标：网络中的所有实例
   来源 IP 范围：0.0.0.0/0
   协议和端口：tcp:8080
   ```

4. **创建并应用**

---

## DigitalOcean Droplet

1. **登录 DigitalOcean**
   - 访问：https://cloud.digitalocean.com/

2. **配置防火墙**
   - 点击 **Networking** > **Firewalls**
   - 选择或创建防火墙
   - 在 **Inbound Rules** 部分添加：
   ```
   Type: Custom
   Protocol: TCP
   Port Range: 8080
   Sources: All IPv4, All IPv6
   ```

---

## 🧪 验证配置

### 1. 获取你的公网 IP

在服务器上执行：
```bash
curl ifconfig.me
```

### 2. 测试本地访问

```bash
curl http://localhost:8080
```

如果返回 HTML 内容，说明服务正常。

### 3. 测试外网访问

在**你的本地电脑**上（不是服务器）：

```bash
# 替换为你的实际公网 IP
curl http://你的公网IP:8080
```

或者直接在浏览器中访问：
```
http://你的公网IP:8080
```

### 4. 检查端口是否开放

使用在线工具：
- https://www.yougetsignal.com/tools/open-ports/
- https://www.portchecktool.com/

---

## 🔒 安全建议

### ⚠️ 警告：开放到公网的风险

将 8080 端口开放到 `0.0.0.0/0` 意味着任何人都可以访问你的应用。

### 推荐方案：限制访问来源

#### 方案 1：仅允许特定 IP
如果你有固定 IP，修改安全组规则：
```
授权对象：你的公网IP/32
```

例如：`203.0.113.1/32`

#### 方案 2：使用 VPN
只开放 VPN IP 段访问：
```
授权对象：10.0.0.0/8
```

#### 方案 3：使用 Nginx + HTTPS + 基本认证

安装 Nginx：
```bash
sudo apt-get update
sudo apt-get install nginx
```

配置反向代理：
```nginx
# /etc/nginx/sites-available/willpower-forge
server {
    listen 80;
    server_name your-domain.com;

    # 基本认证
    auth_basic "Restricted Access";
    auth_basic_user_file /etc/nginx/.htpasswd;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

创建密码文件：
```bash
sudo apt-get install apache2-utils
sudo htpasswd -c /etc/nginx/.htpasswd admin
```

启用配置：
```bash
sudo ln -s /etc/nginx/sites-available/willpower-forge /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

这样你只需要开放 80 (HTTP) 或 443 (HTTPS) 端口，8080 端口保持对外关闭。

#### 方案 4：配置 Let's Encrypt SSL（推荐）

```bash
# 安装 Certbot
sudo apt-get install certbot python3-certbot-nginx

# 获取 SSL 证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

---

## 📊 完整检查清单

- [ ] 服务正在运行：`sudo systemctl status willpower-forge`
- [ ] 服务监听 0.0.0.0:8080：`sudo ss -tulpn | grep 8080`
- [ ] 本地可以访问：`curl http://localhost:8080`
- [ ] 云服务器安全组已配置：开放 8080 端口
- [ ] 外网可以访问：`curl http://公网IP:8080`
- [ ] 浏览器可以访问：`http://公网IP:8080`

---

## 🐛 仍然无法访问？

### 检查监听地址
```bash
sudo ss -tulpn | grep 8080
```

应该看到：
```
tcp   LISTEN 0      4096    *:8080    *:*
```

**不应该是**：
```
tcp   LISTEN 0      4096    127.0.0.1:8080    *:*
```

### 检查服务日志
```bash
sudo journalctl -u willpower-forge -n 50
```

查找 "Listening and serving HTTP on" 这一行。

### 手动测试端口
```bash
# 在服务器上
nc -l 8080

# 在本地电脑上测试
telnet 公网IP 8080
```

### 检查网络连通性
```bash
# 在服务器上
curl ifconfig.me  # 获取公网 IP

# 在本地电脑上
ping 公网IP
traceroute 公网IP
```

---

## 📞 获取帮助

如果仍然无法解决：

1. **提供以下信息**：
   ```bash
   # 运行这些命令并提供输出
   sudo systemctl status willpower-forge
   sudo ss -tulpn | grep 8080
   curl http://localhost:8080
   curl ifconfig.me
   ```

2. **说明你的环境**：
   - 云服务商（阿里云/腾讯云/AWS 等）
   - 操作系统版本
   - 错误信息截图

3. **提交 Issue**：
   - https://github.com/ShengWang1017/willpower/issues

---

**祝配置顺利！** 🚀
