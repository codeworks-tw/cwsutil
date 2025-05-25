# Environment Variables:

| Variable                | module  |        | Type                                        |
| ----------------------- | ------- | ------ | ------------------------------------------- |
| CRYPTO_KEY_HEX          | cwsbase | string | try "openssl rand -hex 32" to generate keys |
| CRYPTO_IV_HEX           | cwsbase | string | try "openssl rand -hex 16" to generate keys |
| ENV                     | cwsbase | string | "test"/"prod"                               |
| IS_LOCAL                | cwsbase | bool   | "true"/"false"/"1"/"0"                      |
| DEBUG                   | cwsbase | bool   | "true"/"false"/"1"/"0"                      |
| LOCALIZATION_LANGUAGE   | cwsbase | string | "en"/"zh_tw"/"zh_cn" (default: en)          |
| S3CacheTTL              | cwsaws  | int    | (default: 10)                               |
| CLOUDWATCHLOG_LOG_GROUP | cwsaws  | string | any aws cloudwatch log group name           |

* *S3CacheTTL*: S3 Object local cache time to live in minutes.

# Release
* 0.1.0 - Apr. 11, 2024
* 0.3.5 - May. 25, 2025

# 使用方法

## ㄧ： 實作 Gin API Handler 完成統一的錯誤處理

1. 實作多語系 json 資料，如下範例，可使用環境變數 LOCALIZATION_LANGUAGE=zh_tw 設定使用中文，預設為 "en" :

```go
import "github.com/codeworks-tw/cwsutil/cwsbase"

const (
    LocalCode_DataNotFound    cwsbase.LocalizationCode = "10000"
)

const LocalizationData string = `{
    "en": {
        "10000": "Data not found - Id: %s",
    },
    "zh_tw": {
        "10000": "資料不存在 - Id: %s",
    }
}`

cwsbase.UpdateLocalizationData([]byte(LocalizationData))  // 更新 cwsbase 多語系目錄
```

2. 使用 cwsutil.WrapHandler 包裝 api handler 來統一 Handle Custom Error cwsutil.CWSError，以達到統一的錯誤 response body，如以下範例:
  
```go
var CHandlerGetData gin.HandlerFunc = cwsutil.WrapHandler(func(ctx *gin.Context) error {
    id := "1"
    var data DataItem
    err := RepositoryData.ToLazyMongoRepository().Get(ctx, cwslazymongo.Eq("_id", id), &data)

    // 處理找不到 id = "1" 的 data
    if err == mongo.ErrNoDocuments {
        return cwsutil.CWSError{
            StatusCode: http.StatusNotFound,
            LocalCode:  LocalCode_DataNotFound,
            EmbeddingStrings: []string{id},
            ActualError: err  // 通常情況不指定此參數，必要時可以使用此參數在 data 欄位回傳真正的錯誤 message
        }
    } else { // 處理其他錯誤
        return err  // 回傳 500 Internal Server Error 並記錄 Error Log
    }

    cwsutil.WriteResponse(ctx, http.StatusOK, cwsbase.LocalCode_OK, data)
    return nil
})
```

   * 有指定 ActualError 的情況， response body 如下:

```json
    // http status: 404 Not Found
    {
        "code": "10000",
        "data": "mongo: no documents in result",
        "message": "Data not found - Id: 1"
    }
```

   * ActualError 為 nil 的情況，response body 如下:

```json
    // http status: 404 Not Found
    {
        "code": "10000",
        "data": null,
        "message": "Data not found - Id: 1"
    }
```