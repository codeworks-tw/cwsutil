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

	"github.com/codeworks-tw/cwsutil/cwsbase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// SetLocalizationData updates the localization data with custom JSON string
// This allows applications to provide their own localized messages
func SetLocalizationData(jsonString string) {
	cwsbase.UpdateLocalizationData([]byte(jsonString))
}

// ParseBody parses the HTTP request body into the provided data structure using Gin's ShouldBind
// Returns a CWSLocalizedErrorResponse error with 400 Bad Request status if parsing fails
func ParseBody(c *gin.Context, data any) error {
	err := c.ShouldBind(data)
	if err != nil {
		return CWSBadRequestError.EmbedActualError(err)
	}
	return nil
}

// ParseQuery parses the HTTP query parameters into the provided data structure using Gin's ShouldBindQuery
// Returns a CWSLocalizedErrorResponse error with 400 Bad Request status if parsing fails
func ParseQuery(c *gin.Context, data any) error {
	err := c.ShouldBindQuery(data)
	if err != nil {
		return CWSBadRequestError.EmbedActualError(err)
	}
	return nil
}

// WrapHandler wraps a function that returns an error into a standard Gin handler
// This provides unified error handling for all HTTP handlers in the application
func WrapHandler(fn func(ctx *gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := fn(ctx)
		if err != nil {
			err = convertDataBaseErrors(err)
			if e, ok := err.(CWSLocalizedErrorResponse); ok {
				if cwsbase.GetEnvironmentInfo().DebugMode {
					log.Println(e)
				}

				if e.actualError != nil {
					WriteResponse(ctx, e.StatusCode, e.LocalCode, e.actualError.Error(), e.embedValues...)
				} else {
					WriteResponse(ctx, e.StatusCode, e.LocalCode, nil, e.embedValues...)
				}
				return
			}
			WriteResponse(ctx, http.StatusInternalServerError, LocalCode_InternalServerError, err.Error())
			panic(err)
		}
	}
}

// WriteResponse writes a standardized JSON response with localized message
// The response format includes code, message, and data fields
func WriteResponse(c *gin.Context, statusCode int, localCode cwsbase.LocalizationCode, data any, localEmbeddingStrs ...any) {
	c.JSON(statusCode, gin.H{
		"code":    localCode,
		"message": cwsbase.GetLocalizationMessage(localCode, localEmbeddingStrs...),
		"data":    data,
	})
}

// WriteResponseWithMongoCursor streams MongoDB cursor data as a JSON response without loading all data into memory
// This is useful for large datasets as it streams data directly from MongoDB to the HTTP response
// Generic type T represents the data type being streamed
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

// convertDataBaseErrors converts common database errors to localized CWSLocalizedErrorResponse errors
// Currently handles GORM "record not found" and MongoDB "no documents" errors
func convertDataBaseErrors(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return CWSNotFoundError.EmbedActualError(err)
	case mongo.ErrNoDocuments:
		return CWSNotFoundError.EmbedActualError(err)
	}
	return err
}
