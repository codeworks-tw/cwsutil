/*
 * File: service.go
 * Created Date: Wednesday, February 14th 2024, 9:56:18 am
 *
 * Last Modified: Sun May 19 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsutil

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codeworks-tw/cwsutil/cwsbase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CWSError struct {
	StatusCode       int
	LocalCode        cwsbase.LocalizationCode
	EmbeddingStrings []any
	ActualError      error
}

func (e CWSError) Error() string {
	r := cwsbase.GetLocalizationMessage(e.LocalCode, e.EmbeddingStrings...)
	if e.ActualError != nil {
		r += " ActualError: " + e.ActualError.Error()
	}
	return r
}

func SetLocalizationData(jsonString string) {
	cwsbase.UpdateLocalizationData([]byte(jsonString))
}

func ParseBody(c *gin.Context, data any) error {
	err := c.ShouldBind(data)
	if err != nil {
		if cwsbase.GetEnvironmentInfo().DebugMode {
			return CWSError{StatusCode: http.StatusBadRequest, LocalCode: cwsbase.LocalCode_BadRequest, ActualError: err}
		} else {
			return CWSError{StatusCode: http.StatusBadRequest, LocalCode: cwsbase.LocalCode_BadRequest, ActualError: nil}
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
	if e, ok := err.(CWSError); ok {
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
