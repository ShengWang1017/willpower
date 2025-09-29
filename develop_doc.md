# AI 驱动开发规约 (Willpower Forge v1.1)

## 1\. 系统概述

- **项目名称:** Willpower Forge
- **目标:** 本文档是一份精确的技术规约，旨在指导一个 AI 代码生成模型完成 Willpower Forge 项目的本地开发版本。AI 必须严格遵守本文档中定义的所有模型、API 接口、业务逻辑和文件结构。
- **核心理念:** 基于《自控力》理论，开发一个 Web 应用，帮助用户通过创建"我要做 (I_WILL)"、"我不要 (I_WONT)"、"我想要 (I_WANT)"三种类型的目标，并每日打卡来训练意志力。

## 2\. 技术栈规约

AI 必须使用以下指定的技术和库版本进行开发。

- **后端 (Backend):**
  - **语言:** Go v1.18 或更高版本
  - **Web 框架:** Gin v1.9.1
  - **数据库:** SQLite (通过 GORM)
  - **ORM:** GORM v1.25.x
  - **ORM 驱动:** GORM SQLite Driver v1.5.x
  - **认证:** github.com/golang-jwt/jwt/v4
  - **密码学:** golang.org/x/crypto/bcrypt
- **前端 (Frontend):**
  - **框架:** Vue.js v3.x (使用 Composition API 和 &lt;script setup&gt;)
  - **构建工具:** Vite v4.x
  - **HTTP 客户端:** Axios v1.x
  - **CSS 框架:** Tailwind CSS v3.x
  - **状态管理:** Pinia v2.x

## 3\. 数据库 Schema 规约

数据库文件名为 willpower.db。AI 需生成 GORM 模型并使用 AutoMigrate 功能在程序启动时创建以下表结构。

### 3.1. 表: users

- **SQL DDL (参考):**  
    CREATE TABLE "users" (  
    "id" INTEGER NOT NULL,  
    "username" TEXT NOT NULL UNIQUE,  
    "password_hash" TEXT NOT NULL,  
    "created_at" DATETIME,  
    "updated_at" DATETIME,  
    PRIMARY KEY ("id" AUTOINCREMENT)  
    );  

- **Go Model (internal/models/user.go):**  
    package models  
    import "time"  
    <br/>type User struct {  
    ID uint \`gorm:"primaryKey" json:"id"\`  
    Username string \`gorm:"unique;not null" json:"username"\`  
    PasswordHash string \`gorm:"not null" json:"-"\` // Omit from JSON responses  
    CreatedAt time.Time \`json:"created_at"\`  
    UpdatedAt time.Time \`json:"updated_at"\`  
    }  

### 3.2. 表: goals

- **SQL DDL (参考):**  
    CREATE TABLE "goals" (  
    "id" INTEGER NOT NULL,  
    "user_id" INTEGER NOT NULL,  
    "type" TEXT NOT NULL,  
    "title" TEXT NOT NULL,  
    "status" TEXT NOT NULL DEFAULT 'active',  
    "created_at" DATETIME,  
    "updated_at" DATETIME,  
    PRIMARY KEY ("id" AUTOINCREMENT),  
    FOREIGN KEY("user_id") REFERENCES "users"("id")  
    );  

- **Go Model (internal/models/goal.go):**  
    package models  
    import "time"  
    <br/>type Goal struct {  
    ID uint \`gorm:"primaryKey" json:"id"\`  
    UserID uint \`gorm:"not null" json:"user_id"\`  
    Type string \`gorm:"not null" json:"type"\` // "I_WILL", "I_WONT", "I_WANT"  
    Title string \`gorm:"not null" json:"title"\`  
    Status string \`gorm:"not null;default:'active'" json:"status"\` // "active", "archived"  
    CreatedAt time.Time \`json:"created_at"\`  
    UpdatedAt time.Time \`json:"updated_at"\`  
    }  

### 3.3. 表: check_ins

- **SQL DDL (参考):**  
    CREATE TABLE "check_ins" (  
    "id" INTEGER NOT NULL,  
    "goal_id" INTEGER NOT NULL,  
    "user_id" INTEGER NOT NULL,  
    "date" TEXT NOT NULL,  
    "status" TEXT NOT NULL,  
    "review_notes" TEXT,  
    "created_at" DATETIME,  
    "updated_at" DATETIME,  
    PRIMARY KEY ("id" AUTOINCREMENT),  
    FOREIGN KEY("goal_id") REFERENCES "goals"("id"),  
    FOREIGN KEY("user_id") REFERENCES "users"("id"),  
    UNIQUE("goal_id", "date")  
    );  

- **Go Model (internal/models/check_in.go):**  
    package models  
    import "time"  
    <br/>type CheckIn struct {  
    ID uint \`gorm:"primaryKey" json:"id"\`  
    GoalID uint \`gorm:"uniqueIndex:idx_goal_date;not null" json:"goal_id"\`  
    UserID uint \`gorm:"not null" json:"user_id"\`  
    Date string \`gorm:"uniqueIndex:idx_goal_date;not null" json:"date"\` // Format: "YYYY-MM-DD"  
    Status string \`gorm:"not null" json:"status"\` // "completed", "failed", "partial"  
    ReviewNotes string \`json:"review_notes"\`  
    CreatedAt time.Time \`json:"created_at"\`  
    UpdatedAt time.Time \`json:"updated_at"\`  
    }  

## 4\. API Endpoint 规约

- **统一前缀:** /api/v1
- **认证方式:** JWT (通过 Authorization: Bearer &lt;token&gt; 请求头传递)
- **通用成功响应结构:**  
    {"code": 0, "message": "Success", "data": { ... } | \[ ... \]}  

- **通用失败响应结构:**  
    {"code": &lt;error_code&gt;, "message": "Error description"}  

### 4.1. 认证接口 (Auth)

**Endpoint:** POST /auth/register

- **描述:** 注册新用户。
- **认证:** 无需
- **请求体 (Request Body):**  
    {  
    "username": "string",  
    "password": "string"  
    }  

- **验证规则:**
  - username: 必需, string, min=3, max=50
  - password: 必需, string, min=8, max=100
- **成功响应 (201 Created):**  
    {"code": 0, "message": "User registered successfully"}  

- **失败响应 (409 Conflict):** username 已存在。  
    {"code": 40901, "message": "Username already exists"}  

- **失败响应 (400 Bad Request):** 输入不合法。  
    {"code": 40001, "message": "Invalid input"}  

**Endpoint:** POST /auth/login

- **描述:** 用户登录并获取 JWT。
- **认证:** 无需
- **请求体 (Request Body):**  
    {  
    "username": "string",  
    "password": "string"  
    }  

- **成功响应 (200 OK):**  
    {"code": 0, "message": "Login successful", "data": {"token": "your_jwt_token_string"}}  

- **失败响应 (401 Unauthorized):** 用户名或密码错误。  
    {"code": 40101, "message": "Invalid username or password"}  

### 4.2. 目标接口 (Goals)

**Endpoint:** POST /goals

- **描述:** 创建一个新目标。
- **认证:** **必需**
- **请求体 (Request Body):**  
    {  
    "type": "string", // Must be one of "I_WILL", "I_WONT", "I_WANT"  
    "title": "string"  
    }  

- **验证规则:**
  - type: 必需, string, in=\["I_WILL", "I_WONT", "I_WANT"\]
  - title: 必需, string, min=1, max=255
- **成功响应 (201 Created):** 返回创建的目标对象。  
    {"code": 0, "message": "Goal created", "data": {"id": 1, "user_id": 1, "type": "I_WILL", ...}}  

**Endpoint:** GET /goals

- **描述:** 获取当前用户的所有目标。
- **认证:** **必需**
- **成功响应 (200 OK):** 返回目标对象数组。  
    {"code": 0, "message": "Success", "data": \[{"id": 1, ...}, {"id": 2, ...}\]}  

### 4.3. 打卡接口 (Check-ins)

**Endpoint:** POST /checkins

- **描述:** 为目标提交今日打卡。如果今日已存在记录，则更新；否则创建新记录。
- **认证:** **必需**
- **请求体 (Request Body):**  
    {  
    "goal_id": "number",  
    "status": "string", // Must be one of "completed", "failed", "partial"  
    "review_notes": "string" // Optional  
    }  

- **验证规则:**
  - goal_id: 必需, number
  - status: 必需, string, in=\["completed", "failed", "partial"\]
- **成功响应 (201 Created / 200 OK):** 返回创建或更新后的打卡记录。  
    {"code": 0, "message": "Check-in recorded", "data": {"id": 101, "goal_id": 1, ...}}  

- **失败响应 (404 Not Found):** goal_id 不存在或不属于当前用户。  
    {"code": 40401, "message": "Goal not found"}  

## 5\. 后端文件结构与逻辑规约

**根目录:** willpower-forge-api

- **go.mod**, **go.sum**
- **main.go:** 程序入口。
  - 初始化配置。
  - 调用 database.Connect() 连接数据库并执行 AutoMigrate。
  - 初始化 Gin 引擎。
  - 调用 routes.SetupRoutes(router) 设置所有路由。
  - 启动 HTTP 服务于 0.0.0.0:8080。
- **internal/database/db.go:**
  - Connect(): 连接 SQLite 数据库，返回 \*gorm.DB 实例。
  - AutoMigrateModels(): 自动迁移所有 models 包下的模型。
- **internal/models/:** 存放 user.go, goal.go, check_in.go 模型文件。
- **internal/handlers/:** 存放 API 处理器。
  - **auth_handler.go:** Register, Login 函数。
  - **goal_handler.go:** CreateGoal, GetGoals 函数。
  - **check_in_handler.go:** CreateOrUpdateCheckIn 函数。
- **internal/services/:** 存放核心业务逻辑。
  - **auth_service.go:**
    - RegisterUser: 哈希密码并创建用户。
    - LoginUser: 验证用户凭据并生成 JWT。
- **internal/middleware/auth.go:**
  - AuthMiddleware(): Gin 中间件。
    - 从 Authorization 头解析 Bearer Token。
    - 验证 Token。
    - 从 Token 中提取 user_id 并存入 gin.Context。
    - 如果 Token 无效，中止请求并返回 401 Unauthorized。
- **internal/routes/router.go:**
  - SetupRoutes():
    - 创建 /api/v1 路由组。
    - 注册 auth 路由 (无需认证)。
    - 创建一个应用了 middleware.AuthMiddleware() 的新路由组。
    - 在认证路由组下注册 goals 和 checkins 路由。

## 6\. 前端文件结构与逻辑规约

**根目录:** willpower-forge-web

- **package.json**, **vite.config.js**
  - vite.config.js 必须配置 server.proxy 将 /api 请求代理到 <http://localhost:8080。>
- **src/**
  - **main.js:** 入口文件，初始化 Vue, Pinia 和 Router。
  - **router/index.js:**
    - 定义路由: /login, /register, / (主看板)。
    - 实现导航守卫: 访问 / 前检查 Pinia store 中是否存在 token，若不存在则重定向到 /login。
  - **store/auth.js:** (Pinia Store)
    - state: token (string), user (object)。
    - actions: login(credentials), register(data), logout()。
    - login action: 调用 API，成功后将 token 存入 localStorage 和 state，然后路由到 /。
  - **services/api.js:**
    - 创建 Axios 实例。
    - 添加请求拦截器: 从 Pinia store 读取 token 并添加到 Authorization 头。
  - **views/:**
    - **LoginPage.vue:** 登录表单。
    - **RegisterPage.vue:** 注册表单。
    - **Dashboard.vue:**
      - 在 onMounted hook 中调用 API 获取 goals 列表。
      - 循环渲染 GoalCard.vue 组件。
  - **components/:**
    - **GoalCard.vue:**
      - props: goal (object)。
      - 显示目标标题和类型。
      - 包含 "完成", "失败", "部分" 三个打卡按钮。
      - 点击按钮时，调用 checkins API 提交打卡记录。# AI 驱动开发规约 (Willpower Forge v1.1)

## 1\. 系统概述

- **项目名称:** Willpower Forge
- **目标:** 本文档是一份精确的技术规约，旨在指导一个 AI 代码生成模型完成 Willpower Forge 项目的本地开发版本。AI 必须严格遵守本文档中定义的所有模型、API 接口、业务逻辑和文件结构。
- **核心理念:** 基于《自控力》理论，开发一个 Web 应用，帮助用户通过创建"我要做 (I_WILL)"、"我不要 (I_WONT)"、"我想要 (I_WANT)"三种类型的目标，并每日打卡来训练意志力。

## 2\. 技术栈规约

AI 必须使用以下指定的技术和库版本进行开发。

- **后端 (Backend):**
  - **语言:** Go v1.18 或更高版本
  - **Web 框架:** Gin v1.9.1
  - **数据库:** SQLite (通过 GORM)
  - **ORM:** GORM v1.25.x
  - **ORM 驱动:** GORM SQLite Driver v1.5.x
  - **认证:** github.com/golang-jwt/jwt/v4
  - **密码学:** golang.org/x/crypto/bcrypt
- **前端 (Frontend):**
  - **框架:** Vue.js v3.x (使用 Composition API 和 &lt;script setup&gt;)
  - **构建工具:** Vite v4.x
  - **HTTP 客户端:** Axios v1.x
  - **CSS 框架:** Tailwind CSS v3.x
  - **状态管理:** Pinia v2.x

## 3\. 数据库 Schema 规约

数据库文件名为 willpower.db。AI 需生成 GORM 模型并使用 AutoMigrate 功能在程序启动时创建以下表结构。

### 3.1. 表: users

- **SQL DDL (参考):**  
    CREATE TABLE "users" (  
    "id" INTEGER NOT NULL,  
    "username" TEXT NOT NULL UNIQUE,  
    "password_hash" TEXT NOT NULL,  
    "created_at" DATETIME,  
    "updated_at" DATETIME,  
    PRIMARY KEY ("id" AUTOINCREMENT)  
    );  

- **Go Model (internal/models/user.go):**  
    package models  
    import "time"  
    <br/>type User struct {  
    ID uint \`gorm:"primaryKey" json:"id"\`  
    Username string \`gorm:"unique;not null" json:"username"\`  
    PasswordHash string \`gorm:"not null" json:"-"\` // Omit from JSON responses  
    CreatedAt time.Time \`json:"created_at"\`  
    UpdatedAt time.Time \`json:"updated_at"\`  
    }  

### 3.2. 表: goals

- **SQL DDL (参考):**  
    CREATE TABLE "goals" (  
    "id" INTEGER NOT NULL,  
    "user_id" INTEGER NOT NULL,  
    "type" TEXT NOT NULL,  
    "title" TEXT NOT NULL,  
    "status" TEXT NOT NULL DEFAULT 'active',  
    "created_at" DATETIME,  
    "updated_at" DATETIME,  
    PRIMARY KEY ("id" AUTOINCREMENT),  
    FOREIGN KEY("user_id") REFERENCES "users"("id")  
    );  

- **Go Model (internal/models/goal.go):**  
    package models  
    import "time"  
    <br/>type Goal struct {  
    ID uint \`gorm:"primaryKey" json:"id"\`  
    UserID uint \`gorm:"not null" json:"user_id"\`  
    Type string \`gorm:"not null" json:"type"\` // "I_WILL", "I_WONT", "I_WANT"  
    Title string \`gorm:"not null" json:"title"\`  
    Status string \`gorm:"not null;default:'active'" json:"status"\` // "active", "archived"  
    CreatedAt time.Time \`json:"created_at"\`  
    UpdatedAt time.Time \`json:"updated_at"\`  
    }  

### 3.3. 表: check_ins

- **SQL DDL (参考):**  
    CREATE TABLE "check_ins" (  
    "id" INTEGER NOT NULL,  
    "goal_id" INTEGER NOT NULL,  
    "user_id" INTEGER NOT NULL,  
    "date" TEXT NOT NULL,  
    "status" TEXT NOT NULL,  
    "review_notes" TEXT,  
    "created_at" DATETIME,  
    "updated_at" DATETIME,  
    PRIMARY KEY ("id" AUTOINCREMENT),  
    FOREIGN KEY("goal_id") REFERENCES "goals"("id"),  
    FOREIGN KEY("user_id") REFERENCES "users"("id"),  
    UNIQUE("goal_id", "date")  
    );  

- **Go Model (internal/models/check_in.go):**  
    package models  
    import "time"  
    <br/>type CheckIn struct {  
    ID uint \`gorm:"primaryKey" json:"id"\`  
    GoalID uint \`gorm:"uniqueIndex:idx_goal_date;not null" json:"goal_id"\`  
    UserID uint \`gorm:"not null" json:"user_id"\`  
    Date string \`gorm:"uniqueIndex:idx_goal_date;not null" json:"date"\` // Format: "YYYY-MM-DD"  
    Status string \`gorm:"not null" json:"status"\` // "completed", "failed", "partial"  
    ReviewNotes string \`json:"review_notes"\`  
    CreatedAt time.Time \`json:"created_at"\`  
    UpdatedAt time.Time \`json:"updated_at"\`  
    }  

## 4\. API Endpoint 规约

- **统一前缀:** /api/v1
- **认证方式:** JWT (通过 Authorization: Bearer &lt;token&gt; 请求头传递)
- **通用成功响应结构:**  
    {"code": 0, "message": "Success", "data": { ... } | \[ ... \]}  

- **通用失败响应结构:**  
    {"code": &lt;error_code&gt;, "message": "Error description"}  

### 4.1. 认证接口 (Auth)

**Endpoint:** POST /auth/register

- **描述:** 注册新用户。
- **认证:** 无需
- **请求体 (Request Body):**  
    {  
    "username": "string",  
    "password": "string"  
    }  

- **验证规则:**
  - username: 必需, string, min=3, max=50
  - password: 必需, string, min=8, max=100
- **成功响应 (201 Created):**  
    {"code": 0, "message": "User registered successfully"}  

- **失败响应 (409 Conflict):** username 已存在。  
    {"code": 40901, "message": "Username already exists"}  

- **失败响应 (400 Bad Request):** 输入不合法。  
    {"code": 40001, "message": "Invalid input"}  

**Endpoint:** POST /auth/login

- **描述:** 用户登录并获取 JWT。
- **认证:** 无需
- **请求体 (Request Body):**  
    {  
    "username": "string",  
    "password": "string"  
    }  

- **成功响应 (200 OK):**  
    {"code": 0, "message": "Login successful", "data": {"token": "your_jwt_token_string"}}  

- **失败响应 (401 Unauthorized):** 用户名或密码错误。  
    {"code": 40101, "message": "Invalid username or password"}  

### 4.2. 目标接口 (Goals)

**Endpoint:** POST /goals

- **描述:** 创建一个新目标。
- **认证:** **必需**
- **请求体 (Request Body):**  
    {  
    "type": "string", // Must be one of "I_WILL", "I_WONT", "I_WANT"  
    "title": "string"  
    }  

- **验证规则:**
  - type: 必需, string, in=\["I_WILL", "I_WONT", "I_WANT"\]
  - title: 必需, string, min=1, max=255
- **成功响应 (201 Created):** 返回创建的目标对象。  
    {"code": 0, "message": "Goal created", "data": {"id": 1, "user_id": 1, "type": "I_WILL", ...}}  

**Endpoint:** GET /goals

- **描述:** 获取当前用户的所有目标。
- **认证:** **必需**
- **成功响应 (200 OK):** 返回目标对象数组。  
    {"code": 0, "message": "Success", "data": \[{"id": 1, ...}, {"id": 2, ...}\]}  

### 4.3. 打卡接口 (Check-ins)

**Endpoint:** POST /checkins

- **描述:** 为目标提交今日打卡。如果今日已存在记录，则更新；否则创建新记录。
- **认证:** **必需**
- **请求体 (Request Body):**  
    {  
    "goal_id": "number",  
    "status": "string", // Must be one of "completed", "failed", "partial"  
    "review_notes": "string" // Optional  
    }  

- **验证规则:**
  - goal_id: 必需, number
  - status: 必需, string, in=\["completed", "failed", "partial"\]
- **成功响应 (201 Created / 200 OK):** 返回创建或更新后的打卡记录。  
    {"code": 0, "message": "Check-in recorded", "data": {"id": 101, "goal_id": 1, ...}}  

- **失败响应 (404 Not Found):** goal_id 不存在或不属于当前用户。  
    {"code": 40401, "message": "Goal not found"}  

## 5\. 后端文件结构与逻辑规约

**根目录:** willpower-forge-api

- **go.mod**, **go.sum**
- **main.go:** 程序入口。
  - 初始化配置。
  - 调用 database.Connect() 连接数据库并执行 AutoMigrate。
  - 初始化 Gin 引擎。
  - 调用 routes.SetupRoutes(router) 设置所有路由。
  - 启动 HTTP 服务于 0.0.0.0:8080。
- **internal/database/db.go:**
  - Connect(): 连接 SQLite 数据库，返回 \*gorm.DB 实例。
  - AutoMigrateModels(): 自动迁移所有 models 包下的模型。
- **internal/models/:** 存放 user.go, goal.go, check_in.go 模型文件。
- **internal/handlers/:** 存放 API 处理器。
  - **auth_handler.go:** Register, Login 函数。
  - **goal_handler.go:** CreateGoal, GetGoals 函数。
  - **check_in_handler.go:** CreateOrUpdateCheckIn 函数。
- **internal/services/:** 存放核心业务逻辑。
  - **auth_service.go:**
    - RegisterUser: 哈希密码并创建用户。
    - LoginUser: 验证用户凭据并生成 JWT。
- **internal/middleware/auth.go:**
  - AuthMiddleware(): Gin 中间件。
    - 从 Authorization 头解析 Bearer Token。
    - 验证 Token。
    - 从 Token 中提取 user_id 并存入 gin.Context。
    - 如果 Token 无效，中止请求并返回 401 Unauthorized。
- **internal/routes/router.go:**
  - SetupRoutes():
    - 创建 /api/v1 路由组。
    - 注册 auth 路由 (无需认证)。
    - 创建一个应用了 middleware.AuthMiddleware() 的新路由组。
    - 在认证路由组下注册 goals 和 checkins 路由。

## 6\. 前端文件结构与逻辑规约

**根目录:** willpower-forge-web

- **package.json**, **vite.config.js**
  - vite.config.js 必须配置 server.proxy 将 /api 请求代理到 <http://localhost:8080。>
- **src/**
  - **main.js:** 入口文件，初始化 Vue, Pinia 和 Router。
  - **router/index.js:**
    - 定义路由: /login, /register, / (主看板)。
    - 实现导航守卫: 访问 / 前检查 Pinia store 中是否存在 token，若不存在则重定向到 /login。
  - **store/auth.js:** (Pinia Store)
    - state: token (string), user (object)。
    - actions: login(credentials), register(data), logout()。
    - login action: 调用 API，成功后将 token 存入 localStorage 和 state，然后路由到 /。
  - **services/api.js:**
    - 创建 Axios 实例。
    - 添加请求拦截器: 从 Pinia store 读取 token 并添加到 Authorization 头。
  - **views/:**
    - **LoginPage.vue:** 登录表单。
    - **RegisterPage.vue:** 注册表单。
    - **Dashboard.vue:**
      - 在 onMounted hook 中调用 API 获取 goals 列表。
      - 循环渲染 GoalCard.vue 组件。
  - **components/:**
    - **GoalCard.vue:**
      - props: goal (object)。
      - 显示目标标题和类型。
      - 包含 "完成", "失败", "部分" 三个打卡按钮。
      - 点击按钮时，调用 checkins API 提交打卡记录。
