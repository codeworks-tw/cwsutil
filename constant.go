package cwsutil

import (
	"net/http"

	"github.com/codeworks-tw/cwsutil/cwsbase"
)

// HTTP status code constants for localization
const (
	// LocalCode_InternalServerError represents HTTP 500 status code
	LocalCode_InternalServerError cwsbase.LocalizationCode = "500"
	// LocalCode_Unauthorized represents HTTP 401 status code
	LocalCode_Unauthorized cwsbase.LocalizationCode = "401"
	// LocalCode_BadRequest represents HTTP 400 status code
	LocalCode_BadRequest cwsbase.LocalizationCode = "400"
	// LocalCode_OK represents HTTP 200 status code
	LocalCode_OK cwsbase.LocalizationCode = "200"
	// LocalCode_Forbidden represents HTTP 403 status code
	LocalCode_Forbidden cwsbase.LocalizationCode = "403"
	// LocalCode_NotFound represents HTTP 404 status code
	LocalCode_NotFound cwsbase.LocalizationCode = "404"
)

// localdata contains default localization data for multiple languages (English, Traditional Chinese, Simplified Chinese)
// This provides standard HTTP error messages in different languages
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

// InitBasicLocalizationData initializes the localization system with default multilingual HTTP status messages
// This function should be called during application startup to load standard error messages
func InitBasicLocalizationData() {
	cwsbase.UpdateLocalizationData([]byte(localdata))
}

// CWSLocalizedErrorResponse represents a localized HTTP response structure that implements the error interface
// It contains status code, localization code, embedded values for message formatting, and actual error details
type CWSLocalizedErrorResponse struct {
	// StatusCode represents the HTTP status code to be returned
	StatusCode int
	// LocalCode is the localization key for retrieving localized messages
	LocalCode cwsbase.LocalizationCode
	// embedValues contains values to be embedded in localized messages (for sprintf formatting)
	embedValues []any
	// actualError stores the underlying error for debugging purposes (only shown in debug mode)
	actualError error
}

// Error implements the error interface for CWSLocalizedErrorResponse
// Returns the localized error message, optionally including actual error details in debug mode
func (r CWSLocalizedErrorResponse) Error() string {
	s := cwsbase.GetLocalizationMessage(r.LocalCode, r.embedValues...)
	if r.actualError != nil {
		s += " ActualError: " + r.actualError.Error()
	}
	return s
}

// EmbedValues sets the values to be embedded in the localized message using sprintf formatting
// Returns the same CWSLocalizedErrorResponse instance for method chaining
func (r *CWSLocalizedErrorResponse) EmbedValues(values ...any) *CWSLocalizedErrorResponse {
	r.embedValues = values
	return r
}

// EmbedActualError sets the actual error to be included in debug mode responses
// The actual error is only shown when DEBUG environment variable is set to true
// Returns the same CWSLocalizedErrorResponse instance for method chaining
func (r *CWSLocalizedErrorResponse) EmbedActualError(err error) *CWSLocalizedErrorResponse {
	if cwsbase.GetEnvironmentInfo().DebugMode {
		// only debug mode can see actual error in response
		r.actualError = err
	}
	return r
}

// CWSInternalServerError represents a pre-configured 500 Internal Server Error response
var CWSInternalServerError = CWSLocalizedErrorResponse{
	StatusCode: http.StatusInternalServerError,
	LocalCode:  LocalCode_InternalServerError,
}

// CWSBadRequestError represents a pre-configured 400 Bad Request error response
var CWSBadRequestError = CWSLocalizedErrorResponse{
	StatusCode: http.StatusBadRequest,
	LocalCode:  LocalCode_BadRequest,
}

// CWSNotFoundError represents a pre-configured 404 Not Found error response
var CWSNotFoundError = CWSLocalizedErrorResponse{
	StatusCode: http.StatusNotFound,
	LocalCode:  LocalCode_NotFound,
}

// CWSForbiddenError represents a pre-configured 403 Forbidden error response
var CWSForbiddenError = CWSLocalizedErrorResponse{
	StatusCode: http.StatusForbidden,
	LocalCode:  LocalCode_Forbidden,
}

// CWSUnauthorizedError represents a pre-configured 401 Unauthorized error response
var CWSUnauthorizedError = CWSLocalizedErrorResponse{
	StatusCode: http.StatusUnauthorized,
	LocalCode:  LocalCode_Unauthorized,
}
