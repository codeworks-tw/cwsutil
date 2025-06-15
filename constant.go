package cwsutil

import (
	"net/http"

	"github.com/codeworks-tw/cwsutil/cwsbase"
)

const (
	LocalCode_InternalServerError cwsbase.LocalizationCode = "500"
	LocalCode_Unauthorized        cwsbase.LocalizationCode = "401"
	LocalCode_BadRequest          cwsbase.LocalizationCode = "400"
	LocalCode_OK                  cwsbase.LocalizationCode = "200"
	LocalCode_Forbidden           cwsbase.LocalizationCode = "403"
	LocalCode_NotFound            cwsbase.LocalizationCode = "404"
)

var localdata string = `{
	"en": {
		"500": "Internal server error",
		"401": "Unauthorized",
		"400": "Bad request",
		"200": "OK",
		"403": "Forbidden",
		"404": "Resource not found"
	},
	"zh_tw": {
		"500": "內部伺服器錯誤",
		"401": "未授權",
		"400": "錯誤的請求",
		"200": "成功",
		"403": "禁止訪問",
		"404": "資源未找到"
	},
	"zh_cn": {
		"500": "内部服务器错误",
		"401": "未授权",
		"400": "错误的请求",
		"200": "成功",
		"403": "禁止访问",
		"404": "资源未找到"
	}
}`

func InitBasicLocalizationData() {
	cwsbase.UpdateLocalizationData([]byte(localdata))
}

var CWSResponseInternalServerError = CWSLocalizedResponse{
	StatusCode: http.StatusInternalServerError,
	LocalCode:  LocalCode_InternalServerError,
}

var CWSResponseBadRequest = CWSLocalizedResponse{
	StatusCode: http.StatusBadRequest,
	LocalCode:  LocalCode_BadRequest,
}

var CWSResponseNotFound = CWSLocalizedResponse{
	StatusCode: http.StatusNotFound,
	LocalCode:  LocalCode_NotFound,
}
