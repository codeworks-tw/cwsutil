/*
 * File: service.go
 * Created Date: Wednesday, February 14th 2024, 9:56:18 am
 *
 * Last Modified: Wed Feb 14 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsutil

import (
	"cwsutil/baseutil"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
)

type CWSError struct {
	StatusCode       int
	LocalCode        baseutil.LocalizationCode
	EmbeddingStrings []any
	ActualError      error
}

func (e *CWSError) Error() string {
	r := baseutil.GetLocalizationMessage(e.LocalCode, e.EmbeddingStrings...)
	if e.ActualError != nil {
		r += " ActualError: " + e.ActualError.Error()
	}
	return r
}

func SetLocalizationData(jsonString string) {
	baseutil.UpdateLocalizationData([]byte(jsonString))
}

func ParseBody(c *gin.Context, data any) error {
	err := c.ShouldBind(data)
	if err != nil {
		if baseutil.GetEnvironmentInfo().DebugMode {
			return &CWSError{StatusCode: http.StatusBadRequest, LocalCode: baseutil.LocalCode_BadRequest, ActualError: err}
		} else {
			return &CWSError{StatusCode: http.StatusBadRequest, LocalCode: baseutil.LocalCode_BadRequest, ActualError: nil}
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
	if e, ok := err.(*CWSError); ok {
		if baseutil.GetEnvironmentInfo().DebugMode {
			log.Println(e)
		}

		if e.ActualError != nil {
			WriteResponse(c, e.StatusCode, e.LocalCode, e.ActualError.Error())
		} else {
			WriteResponse(c, e.StatusCode, e.LocalCode, nil)
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
	m, err := baseutil.StructToMapEscapeEmpty(s)
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

func WriteResponse(c *gin.Context, code int, localCode baseutil.LocalizationCode, data any, strs ...any) {
	c.JSON(code, gin.H{
		"code":    localCode,
		"message": baseutil.GetLocalizationMessage(localCode, strs...),
		"data":    data,
	})
}
