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

	"github.com/codeworks-tw/cwsutil/cwsbase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetLocalizationData updates the localization data with custom JSON string
// This allows applications to provide their own localized messages
func SetLocalizationData(jsonString string) {
	cwsbase.UpdateLocalizationData([]byte(jsonString))
}

// ParseBody parses the HTTP request body into the provided data structure using Gin's ShouldBind
// Supports JSON, XML, YAML, and form data binding based on Content-Type header
// Returns a CWSLocalizedErrorResponse with 400 Bad Request status if parsing fails
func ParseBody(c *gin.Context, data any) error {
	err := c.ShouldBind(data)
	if err != nil {
		return BadRequestErrorResponse.EmbedError(err)
	}
	return nil
}

// ParseQuery parses the HTTP query parameters into the provided data structure using Gin's ShouldBindQuery
// Automatically converts query string parameters to the appropriate struct fields based on struct tags
// Returns a CWSLocalizedErrorResponse with 400 Bad Request status if parsing fails
func ParseQuery(c *gin.Context, data any) error {
	err := c.ShouldBindQuery(data)
	if err != nil {
		return BadRequestErrorResponse.EmbedError(err)
	}
	return nil
}

// WrapHandler wraps a function that returns an error into a standard Gin handler
// This provides unified error handling for all HTTP handlers in the application
// If the error is a CWSLocalizedErrorResponse, it writes the localized response
// Otherwise, it writes a 500 Internal Server Error and panics to ensure the error is logged
func WrapHandler(fn func(ctx *gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := fn(ctx)
		if err != nil {
			if resp, ok := err.(CWSLocalizedErrorResponse); ok {
				resp.WriteResponse(ctx)
				return
			}
			InternalServerErrorResponse.EmbedError(err).WriteResponse(ctx)
			panic(err)
		}
	}
}

// WriteResponse writes a standardized JSON response with localized message
// The response format includes code, message, and data fields
// Parameters:
//   - c: Gin context for writing the HTTP response
//   - statusCode: HTTP status code to return
//   - localCode: Localization code for retrieving the appropriate message
//   - data: Response payload data
//   - localEmbeddingStrs: Optional values to embed in the localized message using sprintf formatting
func WriteResponse(c *gin.Context, statusCode int, localCode cwsbase.LocalizationCode, data any, localEmbeddingStrs ...any) {
	c.JSON(statusCode, gin.H{
		"code":    localCode,
		"message": cwsbase.GetLocalizationMessage(localCode, localEmbeddingStrs...),
		"data":    data,
	})
}

// WriteResponseWithMongoCursor streams MongoDB cursor data as a JSON response without loading all data into memory
// This is useful for large datasets as it streams data directly from MongoDB to the HTTP response
// The response is written incrementally, reducing memory usage for large result sets
// Generic type T represents the data type being streamed from the MongoDB cursor
// Parameters:
//   - c: Gin context for writing the HTTP response
//   - statusCode: HTTP status code to return
//   - localCode: Localization code for retrieving the appropriate message
//   - cursor: MongoDB cursor containing the data to stream
//   - localEmbeddingStrs: Optional values to embed in the localized message using sprintf formatting
//
// Returns an error if cursor iteration or JSON encoding fails
func WriteResponseWithMongoCursor[T any](c *gin.Context, statusCode int, localCode cwsbase.LocalizationCode, cursor *mongo.Cursor, localEmbeddingStrs ...any) error {
	c.Writer.WriteHeader(statusCode)
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	header := `{
		"code": "` + string(localCode) + `",
		"message": "` + cwsbase.GetLocalizationMessage(localCode, localEmbeddingStrs...) + `",
		"data": [`

	bottom := `]}`

	// Track whether we need to add comma separators between JSON array elements
	addComma := false
	c.Writer.Write([]byte(header))
	// Iterate through cursor and stream each document as JSON
	for cursor.Next(c) {
		if addComma {
			c.Writer.Write([]byte(",")) // Add comma separator for JSON array
		}

		if cursor.Err() != nil {
			return cursor.Err()
		}
		// Decode cursor document into the specified type T
		var t T
		err := cursor.Decode(&t)
		if err != nil {
			return err
		}
		// Marshal the decoded document back to JSON for streaming
		d, err := json.Marshal(t)
		if err != nil {
			return err
		}

		// Write the JSON document directly to the response stream
		c.Writer.Write(d)
		addComma = true // Set flag to add comma before next element
	}
	// Close the JSON response and flush the buffer to ensure all data is sent
	c.Writer.Write([]byte(bottom))
	c.Writer.Flush()
	return nil
}
