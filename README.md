# CWS Utilities | CWS 工具庫

A comprehensive Go utility library that simplifies Web API development, database operations, encryption, localization, and cloud services integration.

這是一個為 Go 語言開發的綜合性工具庫，提供了多個模組來簡化 Web API 開發、資料庫操作、加密、多語系支援和雲端服務整合等功能。

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/codeworks-tw/cwsutil)](https://goreportcard.com/report/github.com/codeworks-tw/cwsutil)

## Installation | 安裝

```bash
go get github.com/codeworks-tw/cwsutil
```

## Environment Variables | 環境變數配置

| Variable Name | Module | Type | Description | 說明 |
|---------------|---------|------|-------------|------|
| `CRYPTO_KEY_HEX` | cwsbase | string | Encryption key (generate with `openssl rand -hex 32`) | 加密金鑰 (使用 `openssl rand -hex 32` 產生) |
| `CRYPTO_IV_HEX` | cwsbase | string | Encryption IV (generate with `openssl rand -hex 16`) | 加密向量 (使用 `openssl rand -hex 16` 產生) |
| `ENV` | cwsbase | string | Environment setting: `test`/`prod` | 環境設定: `test`/`prod` |
| `IS_LOCAL` | cwsbase | bool | Local development mode: `true`/`false`/`1`/`0` | 本地開發模式: `true`/`false`/`1`/`0` |
| `DEBUG` | cwsbase | bool | Debug mode: `true`/`false`/`1`/`0` | 除錯模式: `true`/`false`/`1`/`0` |
| `LOCALIZATION_LANGUAGE` | cwsbase | string | Localization setting: `en`/`zh_tw`/`zh_cn` (default: `en`) | 多語系設定: `en`/`zh_tw`/`zh_cn` (預設: `en`) |

## Version History | 版本發佈記錄

- **v0.3.18** - Oct 16, 2025 | 2025年10月16日
- **v0.3.17** - Oct 7, 2025 | 2025年10月7日
- **v0.3.14** - September 28, 2025 | 2025年9月28日
- **v0.3.8** - September 21, 2025 | 2025年9月21日  
- **v0.3.6** - June 15, 2024 | 2024年6月15日
- **v0.3.5** - May 25, 2024 | 2024年5月25日
- **v0.1.0** - April 11, 2024 | 2024年4月11日

## Module Architecture | 模組架構

### Core Modules | 主要模組

- **`cwsutil`** - Main module with HTTP handlers and localized responses | 主模組，提供 HTTP 處理器和本地化回應
- **`cwsbase`** - Core utilities (encryption, localization, environment) | 基礎工具模組（加密、本地化、環境變數）
- **`cwssql`** - SQL database operations with GORM | SQL 資料庫操作模組（使用 GORM）
- **`cwsnosql`** - NoSQL database operations (MongoDB) | NoSQL 資料庫操作模組（MongoDB）
- **`cwsfsm`** - Finite State Machine implementation | 有限狀態機實作模組

## Table of Contents | 目錄

- [Quick Start | 快速開始](#quick-start--快速開始)
- [Core Module (cwsutil) | 主模組](#core-module-cwsutil--主模組)
- [Base Utilities (cwsbase) | 基礎工具](#base-utilities-cwsbase--基礎工具)
- [SQL Database (cwssql) | SQL 資料庫](#sql-database-cwssql--sql-資料庫)
- [NoSQL Database (cwsnosql) | NoSQL 資料庫](#nosql-database-cwsnosql--nosql-資料庫)
- [Finite State Machine (cwsfsm) | 有限狀態機](#finite-state-machine-cwsfsm--有限狀態機)
- [Best Practices | 最佳實踐](#best-practices--最佳實踐)
- [API Response Format | API 回應格式](#api-response-format--api-回應格式)

---

## Quick Start | 快速開始

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/codeworks-tw/cwsutil"
    "github.com/codeworks-tw/cwsutil/cwsbase"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize localization
    // 初始化多語系
    cwsutil.InitBasicLocalizationData()
    
    r := gin.Default()
    
    // Simple handler with unified error handling
    // 使用統一錯誤處理的簡單處理器
    r.GET("/users/:id", cwsutil.WrapHandler(func(ctx *gin.Context) error {
        id := ctx.Param("id")
        
        // Simulate user lookup
        // 模擬用戶查詢
        if id == "" {
        return cwsutil.BadRequestErrorResponse.MessageValues("User ID is required")
        }
        
        // Success response
        // 成功回應
        cwsutil.WriteResponse(ctx, http.StatusOK, cwsutil.LocalCode_OK, map[string]string{
            "id": id,
            "name": "John Doe",
        })
        return nil
    }))
    
    r.Run(":8080")
}
```

---

# Usage Guide | 使用指南

## Core Module (cwsutil) | 主模組

### Unified Error Handling | 統一錯誤處理

The core module provides unified error handling with localized responses:
主模組提供統一的錯誤處理和本地化回應：

```go
import "github.com/codeworks-tw/cwsutil"

// Pre-defined error responses | 預定義錯誤回應
var handler gin.HandlerFunc = cwsutil.WrapHandler(func(ctx *gin.Context) error {
    id := ctx.Param("id")
    
    // Business logic | 業務邏輯
    data, err := findDataById(id)
    if err != nil {
        // Return localized error with embedded values
        // 返回嵌入參數的本地化錯誤
        return cwsutil.NotFoundErrorResponse.MessageValues(id).EmbedError(err)
    }
    
    // Success response | 成功回應
    cwsutil.WriteResponse(ctx, http.StatusOK, cwsutil.LocalCode_OK, data)
    return nil
})
```

### Available Pre-defined Responses | 可用的預定義回應

```go
// Error responses | 錯誤回應
cwsutil.InternalServerErrorResponse  // 500 Internal Server Error
cwsutil.BadRequestErrorResponse      // 400 Bad Request  
cwsutil.UnauthorizedErrorResponse    // 401 Unauthorized
cwsutil.ForbiddenErrorResponse       // 403 Forbidden
cwsutil.NotFoundErrorResponse        // 404 Not Found

// Success response | 成功回應
cwsutil.OKResponse                   // 200 OK
```

### Request Parsing | 請求解析

```go
func handler(ctx *gin.Context) error {
    var requestBody RequestModel
    var queryParams QueryModel
    
    // Parse request body (JSON, XML, YAML, form data)
    // 解析請求主體（JSON、XML、YAML、表單數據）
    if err := cwsutil.ParseBody(ctx, &requestBody); err != nil {
        return err
    }
    
    // Parse query parameters | 解析查詢參數
    if err := cwsutil.ParseQuery(ctx, &queryParams); err != nil {
        return err
    }
    
    // Process business logic... | 處理業務邏輯...
    return nil
}
```

### MongoDB Streaming Response | MongoDB 串流回應

Efficient streaming for large datasets without loading all data into memory:
為大量數據集提供高效串流，不會將所有數據載入內存：

```go
func getUsers(ctx *gin.Context) error {
    cursor, err := getUsersCursor() // Get MongoDB Cursor | 取得 MongoDB Cursor
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    // Stream response without loading all data into memory
    // 串流回應，不會將所有數據載入內存
    return cwsutil.WriteResponseWithMongoCursor[User](ctx, 
        http.StatusOK, cwsutil.LocalCode_OK, cursor)
}
```

---

## Base Utilities (cwsbase) | 基礎工具

### Localization Support | 多語系支援

Comprehensive localization with support for English, Traditional Chinese, and Simplified Chinese:
完整的多語系支援，支援英文、繁體中文和簡體中文：

```go
import "github.com/codeworks-tw/cwsutil/cwsbase"

// Initialize basic localization data | 初始化基本多語系數據
cwsutil.InitBasicLocalizationData()

// Custom localization data | 自定義多語系數據
const LocalizationData = `{
    "en": {
        "10000": "Data not found - Id: %s"
    },
    "zh_tw": {
        "10000": "資料不存在 - Id: %s"
    },
    "zh_cn": {
        "10000": "数据不存在 - Id: %s"
    }
}`

// Update localization data | 更新多語系數據
cwsutil.SetLocalizationData(LocalizationData)

// Get localized message | 取得本地化訊息
message := cwsbase.GetLocalizationMessage("10000", "123")
```

### Environment Variable Management | 環境變數管理

Type-safe environment variable reading with default values:
提供型別安全的環境變數讀取功能，支援預設值：

```go
// Read environment variables with defaults
// 讀取環境變數，支援預設值
port := cwsbase.GetEnv("PORT", 8080)        // int
isDebug := cwsbase.GetEnv("DEBUG", false)  // bool
dbUrl := cwsbase.GetEnv[string]("DB_URL")   // string (required) | 字串（必填）

// Get environment information | 取得環境資訊
envInfo := cwsbase.GetEnvironmentInfo()
fmt.Printf("Env: %s, Debug: %t, Local: %t", envInfo.Env, envInfo.DebugMode, envInfo.IsLocal)
```

### Encryption | 加密功能

AES-CBC encryption and decryption:
提供 AES-CBC 加密/解密功能：

```go
// Encrypt map data | 加密 map 資料
data := map[string]any{
    "userId": 123,
    "role":   "admin",
}
encrypted, err := cwsbase.EncryptMap(data)

// Decrypt back to map data | 解密回 map 資料
decrypted, err := cwsbase.DecryptToMap(encrypted)
```

### Utility Functions | 工具函數

```go
// String processing | 字串處理
camelCase := "HelloWorld"
snakeCase := cwsbase.ToSnakeCase(camelCase) // "hello_world"
capitalized := cwsbase.StringToCapital("hello") // "Hello"

// Time processing | 時間處理
timestamp := cwsbase.GetCurrentTimestamp()        // Unix timestamp | Unix 時間戳
timeStr := cwsbase.GetCurrentTimestampString()    // String format timestamp | 字串格式時間戳
dateTime := cwsbase.IntToDateTime(1234567890)     // Formatted date time | 格式化日期時間

// Struct conversion | 結構轉換
struct2map, err := cwsbase.StructToMap(someStruct)
struct2mapFiltered, err := cwsbase.StructToMapEscapeEmpty(someStruct)

// HTTP requests | HTTP 請求
response, err := cwsbase.SendHttpRequestJson(ctx, "POST", url, jsonBody, headers)
body, err := cwsbase.ReadHttpBody(*response)
```

### Generic Stack Data Structure | 泛型堆疊資料結構

Generic stack implementation:
提供泛型堆疊實作：

```go
// Create stack | 建立堆疊
stack := cwsbase.New[int]()

// Stack operations | 操作堆疊
value := 42
stack.Push(&value)
peek := stack.Peek()    // View top element | 查看頂部元素
popped := stack.Pop()   // Remove top element | 取出頂部元素
length := stack.Len()   // Get length | 取得長度
```

---

## SQL Database (cwssql) | SQL 資料庫

### Database Connection | 資料庫連線

```go
import "github.com/codeworks-tw/cwsutil/cwssql"

// PostgreSQL connection | PostgreSQL 連線
db, err := cwssql.NewPostgresDB("postgres://user:pass@host:port/dbname")

// SQLite connection | SQLite 連線
db, err := cwssql.NewSQLiteDB("./database.sqlite")

// Create session | 建立會話
session := cwssql.NewSession(db)
```

### Base Models | 基礎模型

```go
// Use predefined base models | 使用預定義的基礎模型
type User struct {
    cwssql.BaseIdModel    // Provides UUID primary key | 提供 UUID 主鍵
    cwssql.BaseTimeModel  // Provides CreatedAt, UpdatedAt | 提供 CreatedAt, UpdatedAt
    Name  string          `json:"name"`
    Email string          `json:"email"`
}

// JSON field model | JSON 欄位模型
type Product struct {
    cwssql.BaseIdModel
    cwssql.BaseJsonBModel // PostgreSQL JSONB field | PostgreSQL JSONB 欄位
    Name string            `json:"name"`
}
```

### Repository Pattern | Repository 模式

```go
// Create Repository | 建立 Repository
repo := cwssql.NewRepository[User](ctx, session)

// CRUD operations | CRUD 操作
user := &User{Name: "張三", Email: "zhang@example.com"}

// Create or Update | 新增或更新
err = repo.Upsert(user)

// Query single record | 查詢單一資料
user, err = repo.Get(cwssql.Eq("email", "zhang@example.com"))

// Query multiple records | 查詢多筆資料
users, err = repo.GetAll(cwssql.Like("name", "%張%"))

// Count records | 統計數量
count, err = repo.Count(cwssql.Gte("created_at", yesterday))

// Delete | 刪除
err = repo.Delete(user)
```

### Query Builder | 查詢條件建構器

```go
// Basic conditions | 基本條件
where1 := cwssql.Eq("status", "active")
where2 := cwssql.In("role", "admin", "manager")
where3 := cwssql.Between("age", 18, 65)
where4 := cwssql.Like("name", "%王%")

// Combined conditions | 組合條件
where := cwssql.And(
    cwssql.Eq("status", "active"),
    cwssql.Or(
        cwssql.Eq("role", "admin"),
        cwssql.Gte("experience", 5),
    ),
)

// JSON queries (PostgreSQL) | JSON 查詢 (PostgreSQL)
jsonWhere := cwssql.JSONExtract("metadata", "tags", "important")
jsonbWhere := cwssql.JSONBContains("preferences", `{"theme": "dark"}`)
```

### Transaction Management | 交易處理

```go
// Begin transaction | 開始交易
err = repo.Begin()
if err != nil {
    return err
}

// Execute multiple operations | 執行多個操作
if err = repo.Upsert(user1); err != nil {
    repo.Rollback()
    return err
}

if err = repo.Upsert(user2); err != nil {
    repo.Rollback()
    return err
}

// Commit transaction | 提交交易
if err = repo.Commit(); err != nil {
    return err
}
```

---

## NoSQL Database (cwsnosql) | NoSQL 資料庫

### MongoDB Repository | MongoDB Repository

```go
import "github.com/codeworks-tw/cwsutil/cwsnosql"
import "github.com/codeworks-tw/cwsutil/cwsnosql/cwslazymongo"

// Define primary key structure | 定義主鍵結構
type UserKey struct {
    ID string `bson:"_id"`
}

// Create MongoDB Repository | 建立 MongoDB Repository
repo := &cwsnosql.MongoDBRepository[UserKey]{
    Url:            "mongodb://localhost:27017",
    DbName:         "mydb",
    CollectionName: "users",
}

// Create unique index | 建立唯一索引
err = repo.CreateSimpleUniqueAscendingIndex(ctx)

// CRUD operations | CRUD 操作
key := UserKey{ID: "user123"}
user := User{Name: "李四", Email: "li@example.com"}

// Create or Update | 新增或更新
err = repo.Upsert(ctx, key, user)

// Query | 查詢
var result User
err = repo.Find(ctx, key, &result)

// Delete | 刪除
err = repo.Delete(ctx, key)

// Set operations | 集合操作
_, err = repo.AddValuesToSet(ctx, key, "tags", "golang", "database")
_, err = repo.PullValuesFromSet(ctx, key, "tags", "old-tag")
```

### LazyMongo Query Builder | LazyMongo 查詢建構器

```go
// Get LazyMongo Repository | 取得 LazyMongo Repository
lazyRepo := repo.ToLazyMongoRepository()

// Build query conditions | 建構查詢條件
filter := cwslazymongo.Eq("status", "active").
    Gte("created_at", time.Now().AddDate(0, -1, 0)).
    In("role", "admin", "user")

// Complex queries | 複雜查詢
complexFilter := cwslazymongo.And(
    cwslazymongo.Eq("status", "published"),
    cwslazymongo.Or(
        cwslazymongo.Gt("views", 1000),
        cwslazymongo.Eq("featured", true),
    ),
)

// Execute query | 執行查詢
var users []User
cursor, err := lazyRepo.Select(ctx, filter)
if err == nil {
    defer cursor.Close(ctx)
    for cursor.Next(ctx) {
        var user User
        if err = cursor.Decode(&user); err == nil {
            users = append(users, user)
        }
    }
}
```

### Update Operations | 更新操作

```go
// Build update operations | 建構更新操作
updater := cwslazymongo.Set(map[string]any{
    "last_login": time.Now(),
    "status":     "online",
}).
Inc("login_count", 1).
AddToSet("recent_activities", "login").
Pull("old_notifications", "expired")

// Execute update | 執行更新
result, err := lazyRepo.Update(ctx, filter, updater)
fmt.Printf("Updated %d records", result.ModifiedCount)
```

---

## Finite State Machine (cwsfsm) | 有限狀態機

Powerful state machine implementation for workflows, order processing, approval processes, etc. This module allows you to define discrete steps that can transition between each other based on business logic.

提供強大的狀態機實作，適用於工作流程、訂單處理、審批流程等場景。此模組允許您定義離散的步驟，這些步驟可以根據業務邏輯相互轉換。

### Core Concepts | 核心概念

- **FSMStepName**: Unique identifier for steps | 步驟的唯一識別符
- **FSMStepRegistry**: Registry managing all steps | 管理所有步驟的註冊表
- **FSMStepTransaction**: Transaction object for passing data between steps | 在步驟間傳遞資料的交易物件
- **IFSMStep**: Interface that all steps must implement | 所有步驟必須實現的介面
- **FSMStep**: Functional step implementation | 函數式步驟實作

### Basic Usage | 基本用法

```go
import (
    "context"
    "fmt"
    "github.com/codeworks-tw/cwsutil/cwsfsm"
)

// Define step names | 定義步驟名稱
const (
    StartStep    cwsfsm.FSMStepName = "StartStep"
    ProcessStep  cwsfsm.FSMStepName = "ProcessStep"
    EndStep      cwsfsm.FSMStepName = "EndStep"
)

// Define step implementations | 定義步驟實作
var startStep cwsfsm.FSMStep[int] = func(ctx context.Context, transaction *cwsfsm.FSMStepTransaction[int]) (*cwsfsm.FSMStepTransaction[int], error) {
    fmt.Println("Starting processing, initial value:", transaction.Data)
    transaction.Data = 1
    transaction.NextStep = ProcessStep // Set next step | 設定下一個步驟
    return transaction, nil
}

var processStep cwsfsm.FSMStep[int] = func(ctx context.Context, transaction *cwsfsm.FSMStepTransaction[int]) (*cwsfsm.FSMStepTransaction[int], error) {
    // Increment counter and decide next step | 遞增計數並決定下一步
    if transaction.Data >= 10 {
        transaction.NextStep = EndStep // End processing | 結束處理
        return transaction, nil
    }
    fmt.Println("Processing, current value:", transaction.Data)
    transaction.Data++
    transaction.NextStep = ProcessStep // Continue loop | 繼續循環
    return transaction, nil
}

var endStep cwsfsm.FSMStep[int] = func(ctx context.Context, transaction *cwsfsm.FSMStepTransaction[int]) (*cwsfsm.FSMStepTransaction[int], error) {
    fmt.Println("Processing complete, final value:", transaction.Data)
    // Don't set NextStep, indicating workflow end | 不設定 NextStep，表示工作流程結束
    return transaction, nil
}

// Create step registry | 建立步驟註冊表
stepRegistry := cwsfsm.FSMStepRegistry[int]{
    StartStep:   startStep,
    ProcessStep: processStep,
    EndStep:     endStep,
}

// Execute state machine | 執行狀態機
transaction := &cwsfsm.FSMStepTransaction[int]{
    NextStep: StartStep,
    Data:     0,
}

if err := cwsfsm.RunFSMSteps(context.Background(), stepRegistry, transaction); err != nil {
    fmt.Printf("Execution error: %v\n", err)
}
```

### Dynamic Step Registration | 動態步驟註冊

```go
// Create empty registry | 建立空的註冊表
registry := make(cwsfsm.FSMStepRegistry[string])

// Dynamically add steps | 動態新增步驟
registry.SetFSMStep("step1", func(ctx context.Context, transaction *cwsfsm.FSMStepTransaction[string]) (*cwsfsm.FSMStepTransaction[string], error) {
    transaction.Data += " -> step1"
    transaction.NextStep = "step2"
    return transaction, nil
})

registry.SetFSMStep("step2", func(ctx context.Context, transaction *cwsfsm.FSMStepTransaction[string]) (*cwsfsm.FSMStepTransaction[string], error) {
    transaction.Data += " -> step2"
    return transaction, nil // End workflow | 結束工作流程
})

// Execute dynamically created workflow | 執行動態建立的工作流程
transaction := &cwsfsm.FSMStepTransaction[string]{
    NextStep: "step1",
    Data:     "start",
}

if err := cwsfsm.RunFSMSteps(context.Background(), registry, transaction); err != nil {
    fmt.Printf("Execution error: %v\n", err)
} else {
    fmt.Printf("Result: %s\n", transaction.Data) // Output: "start -> step1 -> step2"
}
```

---

## Best Practices | 最佳實踐

### Project Structure | 專案結構

```
project/
├── main.go
├── config/
│   └── config.go        // Environment variable configuration | 環境變數配置
├── models/
│   ├── user.go          // Data models | 資料模型
│   └── base.go          // Base models | 基礎模型
├── repositories/
│   ├── user_repo.go     // Repository implementation | Repository 實作
│   └── interfaces.go    // Repository interfaces | Repository 介面
├── handlers/
│   └── user_handler.go  // HTTP handlers | HTTP 處理器
├── services/
│   └── user_service.go  // Business logic | 業務邏輯
└── localization/
    └── messages.json    // Localization messages | 多語系訊息
```

### Error Handling Strategy | 錯誤處理策略

```go
// Define project-specific error codes | 定義專案特定的錯誤碼
const (
    LocalCode_UserNotFound   cwsbase.LocalizationCode  = "10001"
    LocalCode_InvalidEmail   cwsbase.LocalizationCode  = "10002"
    LocalCode_DuplicateUser  cwsbase.LocalizationCode  = "10003"
)

// Custom localization data for project-specific errors | 專案特定錯誤的自定義多語系數據
const CustomLocalizationData = `{
    "en": {
        "10001": "User not found - Id: %s",
        "10002": "Invalid email format: %s",
        "10003": "User with email already exists: %s"
    },
    "zh_tw": {
        "10001": "用戶未找到 - Id: %s",
        "10002": "無效的電子郵件格式: %s",
        "10003": "電子郵件已存在的用戶: %s"
    }
}`

// Initialize custom localization data | 初始化自定義多語系數據
cwsutil.SetLocalizationData(CustomLocalizationData)

// Define custom error responses | 定義自定義錯誤回應
var DuplicateUserError = cwsutil.CWSLocalizedErrorResponse{
    StatusCode: http.StatusConflict,
    LocalCode:  LocalCode_DuplicateUser,
}

// Unified error handling in service layer | 在服務層統一處理錯誤
func (s *UserService) CreateUser(user *User) error {
    if !isValidEmail(user.Email) {
        return cwsutil.BadRequestErrorResponse.
            MessageValues(user.Email).
            EmbedError(errors.New("invalid email format"))
    }
    
    if exists, _ := s.repo.UserExists(user.Email); exists {
        return DuplicateUserError.MessageValues(user.Email)
    }
    
    return s.repo.Create(user)
}
```

### Performance Optimization | 性能優化

```go
// Use streaming for large datasets | 對大數據集使用串流
func (h *UserHandler) GetAllUsers(ctx *gin.Context) error {
    cursor, err := h.service.GetUsersCursor()
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    // Stream response to avoid memory issues | 串流回應以避免內存問題
    return cwsutil.WriteResponseWithMongoCursor[User](ctx, 
        http.StatusOK, cwsutil.LocalCode_OK, cursor)
}

// Use transactions for consistency | 使用事務保證一致性
func (s *UserService) TransferCredits(fromUser, toUser *User, amount int) error {
    repo := s.getRepository()
    
    if err := repo.Begin(); err != nil {
        return err
    }
    
    // Execute operations in transaction | 在事務中執行操作
    if err := s.deductCredits(fromUser, amount); err != nil {
        repo.Rollback()
        return err
    }
    
    if err := s.addCredits(toUser, amount); err != nil {
        repo.Rollback() 
        return err
    }
    
    return repo.Commit()
}
```

---

## API Response Format | API 回應格式

All API responses follow a unified format:
所有 API 回應都遵循統一格式：

### Success Response | 成功回應

```json
{
    "code": "200",
    "message": "OK",
    "data": {
        // Actual data | 實際資料
    }
}
```

### Error Response (Production) | 錯誤回應（生產環境）

```json
{
    "code": "10001",
    "message": "User not found - Id: 123",
    "error": null
}
```

### Error Response (Debug Mode) | 錯誤回應（除錯模式）

```json
{
    "code": "10001", 
    "message": "User not found - Id: 123",
    "error": "sql: no rows in result set"
}
```

---

This utility library is designed to provide a clean, secure, and efficient API development experience while maintaining good testability and extensibility. Through unified error handling, localization support, and rich database operation tools, it can significantly improve development efficiency.

這個工具庫的設計理念是提供簡潔、安全、高效的 API 開發體驗，同時保持良好的可測試性和可擴展性。通過統一的錯誤處理、多語系支援和豐富的資料庫操作工具，可以大幅提升開發效率。