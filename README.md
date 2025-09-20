# CWS 工具庫 (CWS Utilities)

這是一個為 Go 語言開發的綜合性工具庫，提供了多個模組來簡化 Web API 開發、資料庫操作、加密、多語系支援等功能。

## 環境變數配置

| 變數名稱              | 所屬模組 | 類型   | 說明                                        |
| --------------------- | -------- | ------ | ------------------------------------------- |
| CRYPTO_KEY_HEX        | cwsbase  | string | 加密金鑰 (使用 "openssl rand -hex 32" 產生) |
| CRYPTO_IV_HEX         | cwsbase  | string | 加密向量 (使用 "openssl rand -hex 16" 產生) |
| ENV                   | cwsbase  | string | 環境設定 "test"/"prod"                      |
| IS_LOCAL              | cwsbase  | bool   | 本地開發模式 "true"/"false"/"1"/"0"         |
| DEBUG                 | cwsbase  | bool   | 除錯模式 "true"/"false"/"1"/"0"             |
| LOCALIZATION_LANGUAGE | cwsbase  | string | 多語系設定 "en"/"zh_tw"/"zh_cn" (預設: en)  |

## 版本發佈記錄
* 0.3.7 - 2025年9月20日
* 0.3.6 - 2024年6月15日
* 0.3.5 - 2024年5月25日
* 0.1.0 - 2024年4月11日

## 模組架構

### 主要模組
- **cwsbase** - 基礎工具模組
- **cwssql** - SQL 資料庫操作模組
- **cwsnosql** - NoSQL 資料庫操作模組
- **cwsfsm** - 有限狀態機模組

---

# 使用指南

## 一、基礎工具模組 (cwsbase)

### 1.1 多語系支援

提供完整的多語系訊息管理功能，支援英文、繁體中文、簡體中文。

#### 基本用法：
```go
import "github.com/codeworks-tw/cwsutil/cwsbase"

// 初始化基本多語系資料
cwsutil.InitBasicLocalizationData()

// 自訂多語系資料
const LocalizationData string = `{
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

// 更新多語系資料
cwsbase.UpdateLocalizationData([]byte(LocalizationData))

// 取得多語系訊息
message := cwsbase.GetLocalizationMessage("10000", "123")
```

### 1.2 環境變數管理

提供型別安全的環境變數讀取功能：

```go
// 讀取環境變數，支援預設值
port := cwsbase.GetEnv("PORT", 8080)        // int
isDebug := cwsbase.GetEnv("DEBUG", false)  // bool
dbUrl := cwsbase.GetEnv[string]("DB_URL")   // string (必填)

// 取得環境資訊
envInfo := cwsbase.GetEnvironmentInfo()
fmt.Printf("環境: %s, 除錯: %t, 本地: %t", envInfo.Env, envInfo.DebugMode, envInfo.IsLocal)
```

### 1.3 加密功能

提供 AES-CBC 加密/解密功能：

```go
// 加密 map 資料
data := map[string]any{
    "userId": 123,
    "role":   "admin",
}
encrypted, err := cwsbase.EncryptMap(data)

// 解密回 map 資料
decrypted, err := cwsbase.DecryptToMap(encrypted)
```

### 1.4 工具函數

```go
// 字串處理
camelCase := "HelloWorld"
snakeCase := cwsbase.ToSnakeCase(camelCase) // "hello_world"
capitalized := cwsbase.StringToCapital("hello") // "Hello"

// 時間處理
timestamp := cwsbase.GetCurrentTimestamp()        // Unix 時間戳
timeStr := cwsbase.GetCurrentTimestampString()    // 字串格式時間戳
dateTime := cwsbase.IntToDateTime(1234567890)     // 格式化日期時間

// 結構轉換
struct2map, err := cwsbase.StructToMap(someStruct)
struct2mapFiltered, err := cwsbase.StructToMapEscapeEmpty(someStruct)

// HTTP 請求
response, err := cwsbase.SendHttpRequestJson(ctx, "POST", url, jsonBody, headers)
body, err := cwsbase.ReadHttpBody(*response)
```

### 1.5 堆疊資料結構

提供泛型堆疊實作：

```go
// 建立堆疊
stack := cwsbase.New[int]()

// 操作堆疊
value := 42
stack.Push(&value)
peek := stack.Peek()    // 查看頂部元素
popped := stack.Pop()   // 取出頂部元素
length := stack.Len()   // 取得長度
```

## 二、Web API 開發 (主模組)

### 2.1 統一錯誤處理

使用 CWSLocalizedErrorResponse 實現統一的錯誤回應格式：

```go
import "github.com/codeworks-tw/cwsutil"

// 定義自訂錯誤回應
var DataNotFoundError = cwsutil.CWSLocalizedErrorResponse{
    StatusCode: http.StatusNotFound,
    LocalCode:  "10000",
}

// 在 Handler 中使用
var handler gin.HandlerFunc = cwsutil.WrapHandler(func(ctx *gin.Context) error {
    id := ctx.Param("id")
    
    // 查詢資料邏輯
    data, err := findDataById(id)
    if err != nil {
        // 回傳自訂錯誤，可嵌入參數和原始錯誤
        return DataNotFoundError.EmbedValues(id).EmbedActualError(err)
    }
    
    // 成功回應
    cwsutil.WriteResponse(ctx, http.StatusOK, cwsutil.LocalCode_OK, data)
    return nil
})
```

### 2.2 請求解析

```go
func handler(ctx *gin.Context) error {
    var requestBody RequestModel
    var queryParams QueryModel
    
    // 解析請求 Body
    if err := cwsutil.ParseBody(ctx, &requestBody); err != nil {
        return err
    }
    
    // 解析查詢參數
    if err := cwsutil.ParseQuery(ctx, &queryParams); err != nil {
        return err
    }
    
    // 處理業務邏輯...
    return nil
}
```

### 2.3 MongoDB 串流回應

針對大量資料提供串流回應功能：

```go
func getUsers(ctx *gin.Context) error {
    cursor, err := getUsersCursor() // 取得 MongoDB Cursor
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    // 串流回應，不會將所有資料載入記憶體
    return cwsutil.WriteResponseWithMongoCursor[User](ctx, 
        http.StatusOK, cwsutil.LocalCode_OK, cursor)
}
```

## 三、SQL 資料庫模組 (cwssql)

### 3.1 資料庫連線

```go
import "github.com/codeworks-tw/cwsutil/cwssql"

// PostgreSQL 連線
db, err := cwssql.NewPostgresDB("postgres://user:pass@host:port/dbname")

// SQLite 連線
db, err := cwssql.NewSQLiteDB("./database.sqlite")

// 建立會話
session := cwssql.NewSession(db)
```

### 3.2 基礎模型

```go
// 使用預定義的基礎模型
type User struct {
    cwssql.BaseIdModel    // 提供 UUID 主鍵
    cwssql.BaseTimeModel  // 提供 CreatedAt, UpdatedAt
    Name  string          `json:"name"`
    Email string          `json:"email"`
}

// JSON 欄位模型
type Product struct {
    cwssql.BaseIdModel
    cwssql.BaseJsonBModel // PostgreSQL JSONB 欄位
    Name string            `json:"name"`
}
```

### 3.3 Repository 模式

```go
// 建立 Repository
repo := cwssql.NewRepository[User](ctx, session)

// CRUD 操作
user := &User{Name: "張三", Email: "zhang@example.com"}

// 新增或更新
err = repo.Upsert(user)

// 查詢單一資料
user, err = repo.Get(cwssql.Eq("email", "zhang@example.com"))

// 查詢多筆資料
users, err = repo.GetAll(cwssql.Like("name", "%張%"))

// 統計數量
count, err = repo.Count(cwssql.Gte("created_at", yesterday))

// 刪除
err = repo.Delete(user)
```

### 3.4 查詢條件建構器

```go
// 基本條件
where1 := cwssql.Eq("status", "active")
where2 := cwssql.In("role", "admin", "manager")
where3 := cwssql.Between("age", 18, 65)
where4 := cwssql.Like("name", "%王%")

// 組合條件
where := cwssql.And(
    cwssql.Eq("status", "active"),
    cwssql.Or(
        cwssql.Eq("role", "admin"),
        cwssql.Gte("experience", 5),
    ),
)

// JSON 查詢 (PostgreSQL)
jsonWhere := cwssql.JSONExtract("metadata", "tags", "important")
jsonbWhere := cwssql.JSONBContains("preferences", `{"theme": "dark"}`)
```

### 3.5 交易處理

```go
// 開始交易
err = repo.Begin()
if err != nil {
    return err
}

// 執行多個操作
if err = repo.Upsert(user1); err != nil {
    repo.Rollback()
    return err
}

if err = repo.Upsert(user2); err != nil {
    repo.Rollback()
    return err
}

// 提交交易
if err = repo.Commit(); err != nil {
    return err
}
```

## 四、NoSQL 資料庫模組 (cwsnosql)

### 4.1 MongoDB Repository

```go
import "github.com/codeworks-tw/cwsutil/cwsnosql"
import "github.com/codeworks-tw/cwsutil/cwsnosql/cwslazymongo"

// 定義主鍵結構
type UserKey struct {
    ID string `bson:"_id"`
}

// 建立 MongoDB Repository
repo := &cwsnosql.MongoDBRepository[UserKey]{
    Url:            "mongodb://localhost:27017",
    DbName:         "mydb",
    CollectionName: "users",
}

// 建立唯一索引
err = repo.CreateSimpleUniqueAscendingIndex(ctx)

// CRUD 操作
key := UserKey{ID: "user123"}
user := User{Name: "李四", Email: "li@example.com"}

// 新增或更新
err = repo.Upsert(ctx, key, user)

// 查詢
var result User
err = repo.Find(ctx, key, &result)

// 刪除
err = repo.Delete(ctx, key)

// 集合操作
_, err = repo.AddValuesToSet(ctx, key, "tags", "golang", "database")
_, err = repo.PullValuesFromSet(ctx, key, "tags", "old-tag")
```

### 4.2 LazyMongo 查詢建構器

```go
// 取得 LazyMongo Repository
lazyRepo := repo.ToLazyMongoRepository()

// 建構查詢條件
filter := cwslazymongo.Eq("status", "active").
    Gte("created_at", time.Now().AddDate(0, -1, 0)).
    In("role", "admin", "user")

// 複雜查詢
complexFilter := cwslazymongo.And(
    cwslazymongo.Eq("status", "published"),
    cwslazymongo.Or(
        cwslazymongo.Gt("views", 1000),
        cwslazymongo.Eq("featured", true),
    ),
)

// 執行查詢
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

### 4.3 更新操作

```go
// 建構更新操作
updater := cwslazymongo.Set(map[string]any{
    "last_login": time.Now(),
    "status":     "online",
}).
Inc("login_count", 1).
AddToSet("recent_activities", "login").
Pull("old_notifications", "expired")

// 執行更新
result, err := lazyRepo.Update(ctx, filter, updater)
fmt.Printf("更新了 %d 筆資料", result.ModifiedCount)
```

## 五、有限狀態機模組 (cwsfsm)

提供強大的狀態機實作，適用於工作流程、訂單處理等場景：

### 5.1 基本用法

```go
import "github.com/codeworks-tw/cwsutil/cwsfsm"

// TestStep implements IFSMStep[int]
var TestStep FSMStep[int] = func(ctx context.Context, transaction *FSMStepTransaction[int]) (*FSMStepTransaction[int], error) {
	// Increment count and continue
	if transaction.Data >= 10 {
		transaction.NextStep = TestEndStep // next to end step
		return transaction, nil
	}
	fmt.Println("Count:", transaction.Data)
	transaction.Data++
	transaction.NextStep = transaction.CurrentStep // loop
	return transaction, nil
}

var TestEndStep FSMStep[int] = func(ctx context.Context, transaction *FSMStepTransaction[int]) (*FSMStepTransaction[int], error) {
	// End count
	fmt.Println("Final:", transaction.Data)
	return nil, nil
}

// Create an action with test steps
if err := RunFSMSetps(context.Background(), &FSMStepTransaction[int]{
		NextStep: TestStep,
		Data:     0,
}); err != nil {
		t.Error(err)
}
```

## 六、預定義錯誤回應

系統提供常用的 HTTP 錯誤回應：

```go
// 直接使用預定義錯誤
return cwsutil.CWSBadRequestError.EmbedValues("invalid parameter")
return cwsutil.CWSNotFoundError.EmbedValues(resourceId)
return cwsutil.CWSUnauthorizedError
return cwsutil.CWSForbiddenError
return cwsutil.CWSInternalServerError.EmbedActualError(err)

// 成功回應
cwsutil.WriteResponse(ctx, cwsutil.CWSOKResponse.StatusCode, cwsutil.CWSOKResponse.LocalCode, data)
```

## 七、最佳實踐建議

### 7.1 專案結構
```
project/
├── main.go
├── config/
│   └── config.go        // 環境變數配置
├── models/
│   ├── user.go          // 資料模型
│   └── base.go          // 基礎模型
├── repositories/
│   ├── user_repo.go     // Repository 實作
│   └── interfaces.go    // Repository 介面
├── handlers/
│   └── user_handler.go  // HTTP 處理器
├── services/
│   └── user_service.go  // 業務邏輯
└── localization/
    └── messages.json    // 多語系訊息
```

### 7.2 錯誤處理策略

```go
// 定義專案特定的錯誤碼
const (
    LocalCode_UserNotFound     = "USER_001"
    LocalCode_InvalidEmail     = "USER_002"
    LocalCode_DuplicateUser    = "USER_003"
)

// 在服務層統一處理錯誤
func (s *UserService) CreateUser(user *User) error {
    if !isValidEmail(user.Email) {
        return cwsutil.CWSBadRequestError.
            EmbedValues(user.Email).
            EmbedActualError(errors.New("invalid email format"))
    }
    
    if exists, _ := s.repo.UserExists(user.Email); exists {
        return customErrors.DuplicateUserError.EmbedValues(user.Email)
    }
    
    return s.repo.Create(user)
}
```

---

## 錯誤回應格式

所有 API 回應都遵循統一格式：

### 成功回應
```json
{
    "code": "200",
    "message": "成功",
    "data": {
        // 實際資料
    }
}
```

### 錯誤回應 (一般模式)
```json
{
    "code": "USER_001",
    "message": "使用者不存在 - Id: 123",
    "data": null
}
```

### 錯誤回應 (除錯模式)
```json
{
    "code": "USER_001",
    "message": "使用者不存在 - Id: 123",
    "data": "sql: no rows in result set"
}
```

---

這個工具庫的設計理念是提供簡潔、安全、高效的 API 開發體驗，同時保持良好的可測試性和可擴展性。通過統一的錯誤處理、多語系支援和豐富的資料庫操作工具，可以大幅提升開發效率。
