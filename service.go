/*
 * File: service.go
 * Created Date: Sunday, May 19th 2024, 2:02:39 pm
 *
 * Last Modified: Mon Jul 22 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsutil

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codeworks-tw/cwsutil/cwsaws"
	"github.com/codeworks-tw/cwsutil/cwsbase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CWSError struct {
	StatusCode       int
	LocalCode        cwsbase.LocalizationCode
	EmbeddingStrings []any
	ActualError      error
}

type CWSLocalizedResponse struct {
	StatusCode       int
	LocalCode        cwsbase.LocalizationCode
	EmbeddingStrings []any
	ActualError      error
}

func (e CWSError) Error() string {
	s := cwsbase.GetLocalizationMessage(e.LocalCode, e.EmbeddingStrings...)
	if e.ActualError != nil {
		s += " ActualError: " + e.ActualError.Error()
	}
	return s
}

func (r CWSLocalizedResponse) Error() string {
	s := cwsbase.GetLocalizationMessage(r.LocalCode, r.EmbeddingStrings...)
	if r.ActualError != nil {
		s += " ActualError: " + r.ActualError.Error()
	}
	return s
}

func (e *CWSLocalizedResponse) EmbedValues(values ...any) CWSLocalizedResponse {
	return CWSLocalizedResponse{
		StatusCode:       e.StatusCode,
		LocalCode:        e.LocalCode,
		EmbeddingStrings: values,
		ActualError:      e.ActualError,
	}
}

func (e *CWSLocalizedResponse) EmbedActualError(err error) CWSLocalizedResponse {
	return CWSLocalizedResponse{
		StatusCode:       e.StatusCode,
		LocalCode:        e.LocalCode,
		EmbeddingStrings: e.EmbeddingStrings,
		ActualError:      err,
	}
}

func SetLocalizationData(jsonString string) {
	cwsbase.UpdateLocalizationData([]byte(jsonString))
}

func ParseBody(c *gin.Context, data any) error {
	err := c.ShouldBind(data)
	if err != nil {
		if cwsbase.GetEnvironmentInfo().DebugMode {
			return CWSResponseBadRequest.EmbedActualError(err)
		} else {
			return CWSResponseBadRequest
		}
	}
	return nil
}

func ParseQuery(c *gin.Context, data any) error {
	err := c.ShouldBindQuery(data)
	if err != nil {
		if cwsbase.GetEnvironmentInfo().DebugMode {
			return CWSResponseBadRequest.EmbedActualError(err)
		} else {
			return CWSResponseBadRequest
		}
	}
	return nil
}

func WrapHandler(fn func(ctx *gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := fn(ctx)
		if err != nil {
			HandleServiceErrors(ctx, err)
		}
	}
}

func HandleServiceErrors(c *gin.Context, err error) {
	err = convertGORMErrors(err)
	if e, ok := err.(CWSLocalizedResponse); ok {
		if cwsbase.GetEnvironmentInfo().DebugMode {
			log.Println(e)
		}

		if e.ActualError != nil {
			WriteResponse(c, e.StatusCode, e.LocalCode, e.ActualError.Error(), e.EmbeddingStrings...)
		} else {
			WriteResponse(c, e.StatusCode, e.LocalCode, nil, e.EmbeddingStrings...)
		}
		return
	}

	logGroup := cwsbase.GetEnv("CLOUDWATCHLOG_LOG_GROUP", "")
	if logGroup != "" {
		proxy := cwsaws.GetCloudWatchLogProxy(logGroup, c)
		e := proxy.SendMessage(err.Error())
		if e != nil {
			log.Println(e)
		}
	}

	WriteResponse(c, http.StatusInternalServerError, LocalCode_InternalServerError, err.Error())
	panic(err)
}

func StructToMap(obj any) (map[string]any, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var data map[string]any
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func StructToAttributeValueMap(s any, modify ...func(key string, val any) any) (map[string]types.AttributeValue, error) {
	m, err := cwsbase.StructToMapEscapeEmpty(s)
	if err != nil {
		return nil, err
	}

	if len(modify) > 0 {
		for k, v := range m {
			m[k] = modify[0](k, v)
		}
	}

	result, err := attributevalue.MarshalMap(m)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func WriteResponse(c *gin.Context, statusCode int, localCode cwsbase.LocalizationCode, data any, localEmbeddingStrs ...any) {
	c.JSON(statusCode, gin.H{
		"code":    localCode,
		"message": cwsbase.GetLocalizationMessage(localCode, localEmbeddingStrs...),
		"data":    data,
	})
}

func WriteResponseWithMongoCursor[T any](c *gin.Context, statusCode int, localCode cwsbase.LocalizationCode, cursor *mongo.Cursor, localEmbeddingStrs ...any) error {
	c.Writer.WriteHeader(statusCode)
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	header := `{
		"code": "` + string(localCode) + `",
		"message": "` + cwsbase.GetLocalizationMessage(localCode, localEmbeddingStrs...) + `",
		"data": [`

	bottom := `]}`

	addComma := false
	c.Writer.Write([]byte(header))
	for cursor.Next(c) {
		if addComma {
			c.Writer.Write([]byte(","))
		}

		if cursor.Err() != nil {
			return cursor.Err()
		}
		var t T
		err := cursor.Decode(&t)
		if err != nil {
			return err
		}
		d, err := json.Marshal(t)
		if err != nil {
			return err
		}

		c.Writer.Write(d)
		addComma = true
	}
	c.Writer.Write([]byte(bottom))
	c.Writer.Flush()
	return nil
}

func IsStringInSlice(val string, ss []string) bool {
	if len(val) == 0 {
		return false
	}
	for _, s := range ss {
		if val == s {
			return true
		}
	}
	return false
}

func convertGORMErrors(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return CWSResponseNotFound
	}
	return err
}
