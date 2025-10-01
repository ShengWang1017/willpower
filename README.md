# Willpower Forge 🎯

<div align="center">

**A beautiful and powerful goal tracking application to forge your willpower**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.3+-4FC08D?logo=vue.js)](https://vuejs.org)

[Features](#-features) • [Quick Start](#-quick-start) • [Screenshots](#-screenshots) • [Development](#-development) • [Building](#-building)

</div>

---

## 📖 Overview

Willpower Forge is a modern goal tracking application that helps you build and maintain productive habits. With an elegant interface featuring dynamic backgrounds, interactive charts, and comprehensive check-in history, staying motivated has never been easier.

The entire application is packaged into a single executable file - no installation, no dependencies, just download and run!

## ✨ Features

- **📊 Visual Progress Tracking** - Beautiful charts and progress rings to visualize your journey
- **✅ Daily Check-ins** - Mark your goals as completed, partial, or failed with review notes
- **📈 Historical Analytics** - View detailed check-in history and trend analysis
- **🎨 Dynamic UI** - Stunning visual effects including particle systems, aurora backgrounds, and shader animations
- **♻️ Recycle Bin** - Recover accidentally deleted goals
- **🌍 Multi-language Support** - Available in English and Chinese
- **🔒 User Authentication** - Secure login and registration system
- **📱 Responsive Design** - Works seamlessly on desktop and mobile devices
- **💾 Standalone Application** - Single executable with embedded frontend, no external dependencies needed

## 🚀 Quick Start

### For Users (Pre-built Executables)

1. **Download** the latest release for your platform:
   - [Windows](https://github.com/ShengWang1017/willpower/releases/latest) - `willpower-forge-windows.exe`
   - [Linux](https://github.com/ShengWang1017/willpower/releases/latest) - `willpower-forge-linux`

2. **Run** the application:

   **Windows:**
   ```cmd
   # Simply double-click willpower-forge-windows.exe
   # Or run from command line:
   willpower-forge-windows.exe
   ```

   **Linux:**
   ```bash
   chmod +x willpower-forge-linux
   ./willpower-forge-linux
   ```

3. **Access** the application:
   - Open your browser and navigate to `http://localhost:8080`
   - Register a new account or login
   - Start tracking your goals!

That's it! No Node.js, Go, or any other runtime required.

## 📸 Screenshots

### Dashboard
Beautiful goal cards with progress visualization and quick check-in actions.

### Goal Detail & History
View comprehensive check-in history with charts showing your progress over time.

### Dynamic Backgrounds
Enjoy stunning visual effects that enhance your user experience.

## 🛠 Development

### Prerequisites

- **Go** 1.18 or higher
- **Node.js** 14+ and npm
- **Git**

### Setup

1. **Clone the repository:**
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

3. **Frontend setup** (in a new terminal):
   ```bash
   cd willpower-forge-web
   npm install
   npm run dev
   ```

4. **Access the development server:**
   - Frontend: `http://localhost:5173`
   - Backend API: `http://localhost:8080`

### Project Structure

```
willpower/
├── willpower-forge-api/          # Go backend (Gin framework)
│   ├── internal/
│   │   ├── database/             # Database connection & models
│   │   ├── handlers/             # HTTP request handlers
│   │   ├── middleware/           # Authentication middleware
│   │   ├── routes/               # API routes
│   │   └── services/             # Business logic
│   ├── web/dist/                 # Embedded frontend files
│   └── main.go                   # Application entry point
├── willpower-forge-web/          # Vue 3 frontend
│   ├── src/
│   │   ├── components/           # Reusable Vue components
│   │   ├── views/                # Page components
│   │   ├── stores/               # Pinia state management
│   │   └── router/               # Vue Router configuration
│   └── dist/                     # Build output
├── build-linux.sh                # Linux build script
├── build-windows.bat             # Windows build script
└── BUILD_README.md               # Detailed build instructions
```

## 🔨 Building

### Building from Source

**Windows:**
```cmd
build-windows.bat
```

**Linux/Mac:**
```bash
chmod +x build-linux.sh
./build-linux.sh
```

The build script will:
1. Install frontend dependencies
2. Build the Vue application
3. Copy frontend files to the backend
4. Compile the Go application with embedded frontend
5. Generate a single executable file

Output files:
- `willpower-forge-windows.exe` (~15MB)
- `willpower-forge-linux` (~18MB)

### Cross-platform Building

Build for other platforms from Linux/Mac:

```bash
# Build for Windows (from Linux/Mac)
cd willpower-forge-api
GOOS=windows GOARCH=amd64 go build -o ../willpower-forge-windows.exe .

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o ../willpower-forge-linux .
```

## 🏗 Technology Stack

### Backend
- **Framework:** Gin (Go web framework)
- **Database:** SQLite with GORM
- **Authentication:** JWT tokens
- **CORS:** gin-contrib/cors

### Frontend
- **Framework:** Vue 3 (Composition API)
- **State Management:** Pinia
- **Routing:** Vue Router
- **UI Framework:** Tailwind CSS
- **Charts:** Chart.js + vue-chartjs
- **Animations:** GSAP, Three.js
- **HTTP Client:** Axios
- **Build Tool:** Vite

### Deployment
- **Embedding:** Go embed package for single-file distribution
- **Build Scripts:** Cross-platform build automation

## 📚 API Documentation

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login

### Goals
- `GET /api/v1/goals` - Get all goals (requires auth)
- `POST /api/v1/goals` - Create new goal (requires auth)
- `GET /api/v1/goals/:id` - Get single goal (requires auth)
- `PUT /api/v1/goals/:id` - Update goal (requires auth)
- `DELETE /api/v1/goals/:id` - Soft delete goal (requires auth)

### Check-ins
- `GET /api/v1/checkins` - Get check-in history by goal_id (requires auth)
- `POST /api/v1/checkins` - Create check-in record (requires auth)

### Recycle Bin
- `GET /api/v1/goals/recycle-bin` - Get deleted goals (requires auth)
- `PUT /api/v1/goals/:id/restore` - Restore deleted goal (requires auth)
- `DELETE /api/v1/goals/:id/permanent` - Permanently delete goal (requires auth)

## 🔧 Configuration

### Database
- Default database file: `willpower.db` (created automatically in the working directory)
- SQLite is used for simplicity and portability

### Server Port
- Default port: `8080`
- To change, modify the port in `willpower-forge-api/main.go:62`

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🐛 Troubleshooting

### Port Already in Use
If port 8080 is already occupied:
```bash
# Find and kill the process using port 8080
# Linux/Mac:
lsof -ti:8080 | xargs kill -9
# Windows:
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

### Database Issues
If you encounter database errors:
```bash
# Backup and remove the database file
mv willpower.db willpower.db.bak
# Restart the application (a fresh database will be created)
```

### Build Errors
Ensure you have the correct versions:
```bash
go version  # Should be 1.18+
node --version  # Should be 14+
npm --version
```

## 📮 Support

If you encounter any issues or have questions:
- Open an [issue](https://github.com/ShengWang1017/willpower/issues)
- Check existing issues for solutions

## 🙏 Acknowledgments

- Built with [Gin](https://github.com/gin-gonic/gin) and [Vue.js](https://vuejs.org)
- UI components inspired by modern design principles
- Icons and visual effects powered by Three.js and GSAP

---

<div align="center">

**Made with ❤️ by [ShengWang1017](https://github.com/ShengWang1017)**

⭐ Star this repo if you find it useful!

</div>
