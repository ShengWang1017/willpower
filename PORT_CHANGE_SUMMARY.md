# 端口已更改为 5173 ✅

## 🎯 快速修改端口步骤

### 当前状态
- ✅ 默认端口：**5173**
- ✅ 支持自定义端口
- ✅ 服务正在运行

### 如何改回 8080 端口

**方法 1：直接修改（最快）**

```bash
# 1. 编辑服务配置
sudo nano /etc/systemd/system/willpower-forge.service

# 2. 找到这一行：
#    Environment="PORT=5173"
# 改成：
#    Environment="PORT=8080"

# 3. 保存退出（Ctrl+X, Y, Enter）

# 4. 重启服务
sudo systemctl daemon-reload
sudo systemctl restart willpower-forge

# 5. 验证
sudo journalctl -u willpower-forge -n 5
```

**方法 2：使用命令（自动化）**

```bash
sudo systemctl stop willpower-forge
sudo sed -i 's/PORT=5173/PORT=8080/g' /etc/systemd/system/willpower-forge.service
sudo systemctl daemon-reload
sudo systemctl start willpower-forge
sudo systemctl status willpower-forge
```

### 验证端口

```bash
# 查看监听端口
sudo ss -tulpn | grep willpower-forge

# 查看服务日志
sudo journalctl -u willpower-forge -n 10

# 测试访问
curl http://localhost:8080  # 如果改成8080
# 或
curl http://localhost:5173  # 默认端口
```

---

## ⚠️ 重要：更新安全组

修改端口后，**必须更新云服务器安全组规则**！

### 阿里云 ECS 快速配置

1. 访问：https://ecs.console.aliyun.com/
2. 找到你的实例 → 安全组 → 配置规则
3. **删除**旧规则（如 8080）
4. **添加**新规则：
   ```
   协议类型：TCP
   端口范围：5173/5173  （或你的自定义端口）
   授权对象：0.0.0.0/0
   ```

---

## 📊 端口对照表

| 用途 | 建议端口 | 说明 |
|------|---------|------|
| 默认 | 5173 | Vite 开发服务器默认端口 |
| 常用 | 8080 | 传统 Web 应用端口 |
| Node | 3000 | Node.js 应用常用 |
| HTTP | 80 | 需要 Nginx 或特权 |
| HTTPS | 443 | 需要 SSL 证书 |

---

## 📚 相关文档

- **详细端口配置**: `PORT_CONFIG.md`
- **外网访问配置**: `EXTERNAL_ACCESS.md`
- **服务管理**: `SYSTEMD_SERVICE.md`
- **快速开始**: `QUICK_START.md`

---

## 🔗 访问地址

根据你的配置：

- **本地访问**: `http://localhost:PORT`
- **外网访问**: `http://36.227.192.178:PORT`

记得替换 `PORT` 为实际端口号（默认 5173）

---

**需要帮助？** 查看 `PORT_CONFIG.md` 获取完整指南！
