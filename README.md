# Willpower Forge 🎯

<div align="center">

**A beautiful and powerful goal tracking application to forge your willpower**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.3+-4FC08D?logo=vue.js)](https://vuejs.org)

[Features](#-features) • [Quick Start](#-quick-start) • [Configuration](#-configuration) • [Development](#-development) • [API](#-api)

</div>

---

## 📖 Overview

Willpower Forge is a modern goal tracking application that helps you build and maintain productive habits. With an elegant interface featuring dynamic backgrounds, interactive charts, and comprehensive check-in history, staying motivated has never been easier.

**Key Highlights:**
- 💾 Single executable file - no installation needed
- 🔒 Secure authentication with JWT
- 📊 Beautiful data visualizations
- 🌐 Multi-language support (EN/中文)
- ⚡ Fast and lightweight (Go + Vue 3)

---

## ✨ Features

- **📊 Visual Progress Tracking** - Charts and progress rings to visualize your journey
- **✅ Daily Check-ins** - Mark goals as completed, partial, or failed with review notes
- **📈 Historical Analytics** - Detailed check-in history and trend analysis
- **🎨 Dynamic UI** - Particle systems, aurora backgrounds, and shader animations
- **♻️ Recycle Bin** - Recover accidentally deleted goals
- **🔒 User Authentication** - Secure login and registration system
- **📱 Responsive Design** - Works on desktop and mobile devices

---

## 🚀 Quick Start

### Option 1: Download Pre-built Executable (Recommended)

1. **Download** the latest release:
   - [Windows](https://github.com/ShengWang1017/willpower/releases/latest) - `willpower-forge-windows.exe`
   - [Linux](https://github.com/ShengWang1017/willpower/releases/latest) - `willpower-forge-linux`

2. **Run** the application:
   ```bash
   # Linux
   chmod +x willpower-forge-linux
   ./willpower-forge-linux

   # Windows
   willpower-forge-windows.exe
   ```

3. **Access**: Open browser at `http://localhost:5173`

### Option 2: Install as Linux System Service

For production deployment with auto-start on boot:

```bash
# 1. Run installation script
sudo ./install-service.sh

# 2. Verify service is running
sudo systemctl status willpower-forge

# 3. Access the application
# Local: http://localhost:5173
# Remote: http://YOUR_SERVER_IP:5173
```

**Service Management:**
```bash
sudo systemctl start willpower-forge    # Start
sudo systemctl stop willpower-forge     # Stop
sudo systemctl restart willpower-forge  # Restart
sudo journalctl -u willpower-forge -f   # View logs
```

**Uninstall:**
```bash
sudo ./uninstall-service.sh
```

---

## ⚙️ Configuration

### Change Listening Port

**Method 1: Environment Variable (Temporary)**
```bash
PORT=8080 ./willpower-forge-linux
```

**Method 2: Systemd Service (Permanent)**
```bash
# Edit service file
sudo nano /etc/systemd/system/willpower-forge.service

# Change this line:
Environment="PORT=5173"    # Change to your desired port

# Restart service
sudo systemctl daemon-reload
sudo systemctl restart willpower-forge
```

**Default Port:** 5173

### Enable External Access

If running on a cloud server (Aliyun, AWS, etc.), you need to configure security group:

**Aliyun ECS:**
1. Login to https://ecs.console.aliyun.com/
2. Navigate to: Instance → Security Groups → Configure Rules
3. Add Inbound Rule:
   - Protocol: TCP
   - Port Range: 5173/5173 (or your custom port)
   - Authorization Object: 0.0.0.0/0
   - Description: Willpower Forge

**Other Cloud Providers:**
- AWS EC2: Security Groups → Inbound Rules → Add TCP 5173
- Tencent Cloud: Security Groups → Add Rule
- Google Cloud: VPC Network → Firewall Rules

**Verify External Access:**
```bash
# Get your public IP
curl ifconfig.me

# Test access from your local computer
curl http://YOUR_PUBLIC_IP:5173
```

### Database Location

By default, `willpower.db` is created in the working directory:
- Direct run: Same directory as executable
- Service: `/var/lib/willpower-forge/willpower.db`

---

## 🛠 Development

### Prerequisites

- Go 1.18+
- Node.js 14+
- npm

### Setup

1. **Clone repository:**
   ```bash
   git clone https://github.com/ShengWang1017/willpower.git
   cd willpower
   ```

2. **Backend setup:**
   ```bash
   cd willpower-forge-api
   go mod download
   go run main.go
   ```

3. **Frontend setup** (new terminal):
   ```bash
   cd willpower-forge-web
   npm install
   npm run dev
   ```

4. **Access:**
   - Frontend: `http://localhost:5173` (Vite dev server)
   - Backend API: `http://localhost:5173/api/v1`

### Project Structure

```
willpower/
├── willpower-forge-api/       # Go backend (Gin)
│   ├── internal/
│   │   ├── database/          # DB models & connection
│   │   ├── handlers/          # HTTP handlers
│   │   ├── middleware/        # Auth middleware
│   │   ├── routes/            # API routes
│   │   └── services/          # Business logic
│   ├── web/dist/              # Embedded frontend
│   └── main.go
├── willpower-forge-web/       # Vue 3 frontend
│   ├── src/
│   │   ├── components/        # Vue components
│   │   ├── views/             # Page views
│   │   ├── stores/            # Pinia stores
│   │   └── router/            # Vue Router
│   └── dist/                  # Build output
├── install-service.sh         # Service installer
├── uninstall-service.sh       # Service uninstaller
├── build-linux.sh             # Linux build script
├── build-windows.bat          # Windows build script
└── BUILD_README.md            # Build documentation
```

### Building

See [BUILD_README.md](BUILD_README.md) for detailed build instructions.

**Quick Build:**
```bash
# Linux/Mac
./build-linux.sh

# Windows
build-windows.bat
```

### Technology Stack

**Backend:**
- Framework: Gin (Go)
- Database: SQLite + GORM
- Auth: JWT tokens
- Embedding: go:embed for frontend

**Frontend:**
- Framework: Vue 3 (Composition API)
- State: Pinia
- Router: Vue Router
- UI: Tailwind CSS
- Charts: Chart.js
- Animations: GSAP, Three.js
- Build: Vite

---

## 📚 API Documentation

Base URL: `http://localhost:5173/api/v1`

### Authentication

#### Register
```http
POST /auth/register
Content-Type: application/json

{
  "username": "user",
  "password": "password"
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "username": "user",
  "password": "password"
}

Response: { "token": "jwt_token" }
```

### Goals (Requires Authentication)

Add header: `Authorization: Bearer <token>`

#### List Goals
```http
GET /goals
```

#### Create Goal
```http
POST /goals
Content-Type: application/json

{
  "title": "Goal Title",
  "description": "Goal Description"
}
```

#### Get Goal
```http
GET /goals/:id
```

#### Update Goal
```http
PUT /goals/:id
Content-Type: application/json

{
  "title": "Updated Title",
  "description": "Updated Description"
}
```

#### Delete Goal (Soft)
```http
DELETE /goals/:id
```

#### Restore Goal
```http
PUT /goals/:id/restore
```

#### Permanent Delete
```http
DELETE /goals/:id/permanent
```

#### Get Recycle Bin
```http
GET /goals/recycle-bin
```

### Check-ins (Requires Authentication)

#### Create Check-in
```http
POST /checkins
Content-Type: application/json

{
  "goal_id": 1,
  "status": "completed",  // "completed" | "partial" | "failed"
  "review": "Daily review notes"
}
```

#### Get Check-in History
```http
GET /checkins?goal_id=1
```

#### Get Summary
```http
GET /checkins/summary
```

---

## 🐛 Troubleshooting

### Port Already in Use

```bash
# Find process using port
sudo lsof -i:5173
# or
sudo ss -tulpn | grep 5173

# Kill process
sudo kill -9 <PID>
```

### Cannot Access from External Network

1. **Verify service is running:**
   ```bash
   sudo systemctl status willpower-forge
   ```

2. **Check listening address:**
   ```bash
   sudo ss -tulpn | grep 5173
   # Should show: *:5173 (not 127.0.0.1:5173)
   ```

3. **Configure cloud security group** (see [Configuration](#-configuration) section)

4. **Check firewall:**
   ```bash
   # Ubuntu/Debian
   sudo ufw allow 5173/tcp

   # CentOS/RHEL
   sudo firewall-cmd --add-port=5173/tcp --permanent
   sudo firewall-cmd --reload
   ```

### Service Won't Start

```bash
# Check detailed logs
sudo journalctl -u willpower-forge -xe

# Common issues:
# - Port already in use
# - Permission denied (check file permissions)
# - Database locked (stop other instances)
```

### Build Errors

```bash
# Verify versions
go version   # Should be 1.18+
node --version  # Should be 14+

# Clean and rebuild
cd willpower-forge-web
rm -rf node_modules dist
npm install
npm run build

cd ../willpower-forge-api
go clean
go build
```

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License.

---

## 🙏 Acknowledgments

- Built with [Gin](https://github.com/gin-gonic/gin) and [Vue.js](https://vuejs.org)
- UI components inspired by modern design principles
- Icons and visual effects powered by Three.js and GSAP

---

## 📮 Support

- 🐛 Issues: https://github.com/ShengWang1017/willpower/issues
- 📧 Email: (Add your email if needed)
- 📖 Documentation: [BUILD_README.md](BUILD_README.md)

---

<div align="center">

**Made with ❤️ for better habits**

⭐ Star this repo if you find it useful!

</div>
