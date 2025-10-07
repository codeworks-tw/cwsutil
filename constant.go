package cwsutil

import (
	"log"
	"net/http"

	"github.com/codeworks-tw/cwsutil/cwsbase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// HTTP status code constants for localization
const (
	// LocalCode_InternalServerError represents HTTP 500 status code
	LocalCode_InternalServerError cwsbase.LocalizationCode = "500"
	// LocalCode_BadRequest represents HTTP 400 status code
	LocalCode_BadRequest cwsbase.LocalizationCode = "400"
	// LocalCode_Unauthorized represents HTTP 401 status code
	LocalCode_Unauthorized cwsbase.LocalizationCode = "401"
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
		"400": "Bad request",
		"401": "Unauthorized",
		"200": "OK",
		"403": "Forbidden",
		"404": "Resource not found"
	},
	"zh_tw": {
		"500": "內部伺服器錯誤",
		"400": "請求錯誤",
		"401": "未授權",
		"200": "成功",
		"403": "禁止訪問",
		"404": "資源未找到"
	},
	"zh_cn": {
		"500": "内部服务器错误",
		"400": "请求错误",
		"401": "未授权",
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

// CWSLocalizedResponse represents a localized HTTP response structure
// It contains status code, localization code, embedded values for message formatting, and response data
type CWSLocalizedResponse struct {
	// StatusCode represents the HTTP status code to be returned
	StatusCode int
	// LocalCode is the localization key for retrieving localized messages
	LocalCode cwsbase.LocalizationCode
	// embedValues contains values to be embedded in localized messages (for sprintf formatting)
	embedValues []any
	// data contains the response payload data
	data any
}

// ToMessage returns the localized message for this response
func (r CWSLocalizedResponse) ToMessage() string {
	return cwsbase.GetLocalizationMessage(r.LocalCode, r.embedValues...)
}

// WriteResponse writes the HTTP response to the gin context in JSON format
func (r CWSLocalizedResponse) WriteResponse(ctx *gin.Context) {
	ctx.JSON(r.StatusCode, gin.H{
		"code":    r.LocalCode,
		"message": r.ToMessage(),
		"data":    r.data,
	})
}

// MessageValues sets the values to be embedded in the localized message using sprintf formatting
// Returns the same CWSLocalizedResponse instance for method chaining
func (r CWSLocalizedResponse) MessageValues(values ...any) CWSLocalizedResponse {
	r.embedValues = values
	return r
}

// ResponseData sets the data to be included in the response
// Returns the same CWSLocalizedResponse instance for method chaining
func (r CWSLocalizedResponse) ResponseData(data any) CWSLocalizedResponse {
	r.data = data
	return r
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
	// err stores the underlying error for debugging purposes (only shown in debug mode)
	err error
}

// Error implements the error interface for CWSLocalizedErrorResponse
// Returns the localized error message, optionally including actual error details in debug mode
func (r CWSLocalizedErrorResponse) Error() string {
	return r.ToMessage()
}

// ToMessage returns the localized message for this error response
func (r CWSLocalizedErrorResponse) ToMessage() string {
	return cwsbase.GetLocalizationMessage(r.LocalCode, r.embedValues...)
}

// WriteResponse writes the HTTP error response to the gin context in JSON format
// Automatically handles common database errors and debug mode error details
func (r CWSLocalizedErrorResponse) WriteResponse(ctx *gin.Context) {
	if r.err != nil {
		if r.err == gorm.ErrRecordNotFound || r.err == mongo.ErrNoDocuments {
			r.StatusCode = http.StatusNotFound
			r.LocalCode = LocalCode_NotFound
		}
		log.Println(r.err)
		if r.StatusCode < 500 {
			ctx.JSON(r.StatusCode, gin.H{
				"code":    r.LocalCode,
				"message": r.ToMessage(),
				"error":   r.err.Error(),
			})
			return
		} else if r.StatusCode >= 500 && cwsbase.GetEnvironmentInfo().DebugMode {
			ctx.JSON(r.StatusCode, gin.H{
				"code":    r.LocalCode,
				"message": r.ToMessage(),
				"error":   r.err.Error(),
			})
			return
		}
	}
	ctx.JSON(r.StatusCode, gin.H{
		"code":    r.LocalCode,
		"message": r.ToMessage(),
		"error":   nil,
	})
}

// MessageValues sets the values to be embedded in the localized error message using sprintf formatting
// Returns the same CWSLocalizedErrorResponse instance for method chaining
func (r CWSLocalizedErrorResponse) MessageValues(values ...any) CWSLocalizedErrorResponse {
	r.embedValues = values
	return r
}

// EmbedError sets the actual error to be included in debug mode responses
// The actual error is only shown when DEBUG environment variable is set to true
// Returns the same CWSLocalizedErrorResponse instance for method chaining
func (r CWSLocalizedErrorResponse) EmbedError(err error) CWSLocalizedErrorResponse {
	r.err = err
	return r
}

// InternalServerErrorResponse represents a pre-configured 500 Internal Server Error response
var InternalServerErrorResponse = CWSLocalizedErrorResponse{
	StatusCode: http.StatusInternalServerError,
	LocalCode:  LocalCode_InternalServerError,
}

// BadRequestErrorResponse represents a pre-configured 400 Bad Request error response
var BadRequestErrorResponse = CWSLocalizedErrorResponse{
	StatusCode: http.StatusBadRequest,
	LocalCode:  LocalCode_BadRequest,
}

// NotFoundErrorResponse represents a pre-configured 404 Not Found error response
var NotFoundErrorResponse = CWSLocalizedErrorResponse{
	StatusCode: http.StatusNotFound,
	LocalCode:  LocalCode_NotFound,
}

// ForbiddenErrorResponse represents a pre-configured 403 Forbidden error response
var ForbiddenErrorResponse = CWSLocalizedErrorResponse{
	StatusCode: http.StatusForbidden,
	LocalCode:  LocalCode_Forbidden,
}

// UnauthorizedErrorResponse represents a pre-configured 401 Unauthorized error response
var UnauthorizedErrorResponse = CWSLocalizedErrorResponse{
	StatusCode: http.StatusUnauthorized,
	LocalCode:  LocalCode_Unauthorized,
}

// OKResponse represents a pre-configured 200 OK success response
var OKResponse = CWSLocalizedResponse{
	StatusCode: http.StatusOK,
	LocalCode:  LocalCode_OK,
}
