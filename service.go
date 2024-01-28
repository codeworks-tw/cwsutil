/*
 * File: service.go
 * Created Date: Saturday, January 27th 2024, 9:54:26 am
 *
 * Last Modified: Sat Jan 27 2024
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
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
)

type IService interface {
	SetAPIs(ginEngine *gin.Engine)
	SetLocalizationData()
}

func InitializeService(ginEngine *gin.Engine, service IService) {
	service.SetAPIs(ginEngine)
	service.SetLocalizationData()
}

type ServiceError struct {
	StatusCode       int
	LocalCode        baseutil.LocalizationCode
	EmbeddingStrings []any
	ActualError      error
}

func (e *ServiceError) Error() string {
	r := baseutil.GetLocalizationMessage(e.LocalCode, e.EmbeddingStrings...)
	if e.ActualError != nil {
		r += " ActualError: " + e.ActualError.Error()
	}
	return r
}

type BaseHandler struct{}

func (h *BaseHandler) StructToMap(obj any) (map[string]any, error) {
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

func (h *BaseHandler) SetLocalizationData(jsonString string) {
	baseutil.UpdateLocalizationData([]byte(jsonString))
}

func (h *BaseHandler) HandleServiceErrors(c *gin.Context, err error) {
	if e, ok := err.(*ServiceError); ok {
		if baseutil.GetEnvironmentInfo().DebugMode {
			log.Println(e)
		}
		if e.ActualError != nil {
			h.WriteResponse(c, e.StatusCode, e.LocalCode, e.ActualError.Error())
		} else {
			h.WriteResponse(c, e.StatusCode, e.LocalCode, nil)
		}
		return
	}
	panic(err)
}

func (h *BaseHandler) WriteResponse(c *gin.Context, code int, localCode baseutil.LocalizationCode, data any, strs ...any) {
	WriteResponse(c, code, localCode, data, strs...)
}

func (h *BaseHandler) ParseBody(c *gin.Context, data any) error {
	err := c.ShouldBind(data)
	if err != nil {
		if baseutil.GetEnvironmentInfo().DebugMode {
			return &ServiceError{StatusCode: http.StatusBadRequest, LocalCode: baseutil.LocalCode_BadRequest, ActualError: err}
		} else {
			return &ServiceError{StatusCode: http.StatusBadRequest, LocalCode: baseutil.LocalCode_BadRequest, ActualError: nil}
		}
	}
	return nil
}

func WrapHandle(fn func(c *gin.Context, handler *BaseHandler) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := BaseHandler{}
		err := fn(c, &h)
		if err != nil {
			h.HandleServiceErrors(c, err)
		}
	}
}

func HandleServiceErrors(c *gin.Context, err error) {
	if e, ok := err.(*ServiceError); ok {
		if baseutil.GetEnvironmentInfo().DebugMode {
			log.Println(e)
		}
		WriteResponse(c, e.StatusCode, e.LocalCode, nil)
		return
	}
	panic(err)
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
	if strings.Contains(c.FullPath(), "/v1") {
		m := map[string]any{}
		if data != nil {
			if s, ok := data.(string); ok {
				m["data"] = s
			} else {
				j, _ := json.Marshal(data)
				json.Unmarshal(j, &m)
			}
		} else {
			m = make(map[string]any)
		}
		m["error_code"] = localCode
		m["message"] = baseutil.GetLocalizationMessage(localCode, strs...)

		c.JSON(code, m)
	} else {
		c.JSON(code, gin.H{
			"code":    localCode,
			"message": baseutil.GetLocalizationMessage(localCode, strs...),
			"data":    data,
		})
	}
}
