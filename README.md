# Environment Variables:

| Variable              | Type   | Value                                       |
| --------------------- | ------ | ------------------------------------------- |
| CRYPTO_KEY_HEX        | string | try "openssl rand -hex 32" to generate keys |
| CRYPTO_IV_HEX         | string | try "openssl rand -hex 16" to generate keys |
| ENV                   | string | "test"/"prod"                               |
| IS_LOCAL              | bool   | "true"/"false"/"1"/"0"                      |
| DEBUG                 | bool   | "true"/"false"/"1"/"0"                      |
| LOCALIZATION_LANGUAGE | string | "en"/"zh_tw"/"zh_cn" (default: en)          |
| S3CacheTTL            | int    | (default: 10)                               |

* *S3CacheTTL*: S3 Object local cache time to live in minutes.

# Release
* 0.1.0 - Apr. 11, 2024

# 使用方法

## ㄧ： 實作 API 完成統一規範的錯誤處理

1. 實作多語系 json 資料，如下範例，可使用環境變數 LOCALIZATION_LANGUAGE=zh_tw 設定中文，預設為 "en" :

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

2. 使用 cwsutil.WrapHandler 包裝 api handler 來統一 Handle Custom Error cwsutil.CWSError，以達到統一規範回傳錯誤情況的 response body，如以下範例:
  
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