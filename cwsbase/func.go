/*
 * File: func.go
 * Created Date: Thursday, April 11th 2024, 10:31:37 am
 *
 * Last Modified: Tue Jun 04 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsbase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// EnvironmentInfo holds application environment configuration
type EnvironmentInfo struct {
	// Env represents the current environment (e.g., "test", "prod")
	Env       string
	// DebugMode indicates if debug mode is enabled
	DebugMode bool
	// IsLocal indicates if the application is running in local development mode
	IsLocal   bool
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase converts CamelCase or PascalCase strings to snake_case
// Example: "HelloWorld" becomes "hello_world"
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// GetEnvironmentInfo retrieves current environment configuration from environment variables
// Reads ENV, DEBUG, and IS_LOCAL environment variables
func GetEnvironmentInfo() EnvironmentInfo {
	return EnvironmentInfo{
		Env:       GetEnv[string]("ENV"),
		DebugMode: GetEnv("DEBUG", false),
		IsLocal:   GetEnv("IS_LOCAL", false),
	}
}

// GetEnv retrieves an environment variable and converts it to the specified type T
// Supports bool, int, int32, int64, float32, float64, and string types
// If the environment variable is not set and no default value is provided, the program will exit with a fatal error
func GetEnv[T any](key string, defaultVal ...T) T {
	p := os.Getenv(key)
	if p == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		log.Fatalln("Missing required environment variable: " + key)
	}
	var t T
	switch any(t).(type) {
	case bool:
		t, err := strconv.ParseBool(p)
		if err != nil {
			log.Fatalln("Invalid environment bool variable: " + key)
		}
		return any(t).(T)
	case int:
		t, err := strconv.ParseInt(p, 10, 32)
		if err != nil {
			log.Fatalln("Invalid environment int variable: " + key)
		}
		return any(int(t)).(T)
	case int32:
		t, err := strconv.ParseInt(p, 10, 32)
		if err != nil {
			log.Fatalln("Invalid environment int variable: " + key)
		}
		return any(int32(t)).(T)
	case int64:
		t, err := strconv.ParseInt(p, 10, 32)
		if err != nil {
			log.Fatalln("Invalid environment int variable: " + key)
		}
		return any(t).(T)
	case float32:
		t, err := strconv.ParseFloat(p, 32)
		if err != nil {
			log.Fatalln("Invalid environment float variable: " + key)
		}
		return any(float32(t)).(T)
	case float64:
		t, err := strconv.ParseFloat(p, 32)
		if err != nil {
			log.Fatalln("Invalid environment float variable: " + key)
		}
		return any(t).(T)
	default:
		return any(p).(T)
	}
}

// IntToDateTime converts a Unix timestamp to a formatted date-time string
// Returns the time in "2006-01-02 15:04:05" format
func IntToDateTime(unixTime int64) string {
	return time.Unix(unixTime, 0).Format("2006-01-02 15:04:05")
}

// GetCurrentTimestampString returns the current UTC timestamp as a string
func GetCurrentTimestampString() string {
	return strconv.FormatInt(time.Now().UTC().Unix(), 10)
}

// GetCurrentTimestamp returns the current UTC timestamp as an int64
func GetCurrentTimestamp() int64 {
	return time.Now().UTC().Unix()
}

// StructToMapEscapeEmpty converts a struct to a map[string]any while excluding empty values
// Empty values include nil, empty string, and zero integer values
// Returns an error if the input is not a struct
func StructToMapEscapeEmpty(obj any) (map[string]any, error) {
	result := map[string]any{}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expect struct, but got %s", v.Kind())
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fi := t.Field(i)
		switch v.Field(i).Interface() {
		case nil:
		case "":
		case 0:
			continue
		default:
			result[fi.Name] = v.Field(i).Interface()
		}
	}

	return result, nil
}

// MaxInt32 returns the maximum of two int32 values
func MaxInt32(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// AvoidNilString safely dereferences a string pointer, returning empty string if nil
func AvoidNilString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ConvertArrayAnyToArrayString converts a slice of any to a slice of strings
// Assumes all elements in the input slice can be cast to string
func ConvertArrayAnyToArrayString(arr []any) []string {
	var result []string
	for _, v := range arr {
		result = append(result, v.(string))
	}
	return result
}

// ReadHttpBody reads and parses an HTTP response body as JSON into a map[string]any
// Automatically closes the response body when finished
func ReadHttpBody(response http.Response) (map[string]any, error) {
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]any
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SendHttpRequestJson sends an HTTP request with JSON body and custom headers
// The jsonBody is marshaled to JSON and sent as the request body
func SendHttpRequestJson(c context.Context, method string, url string, jsonBody map[string]any, header map[string]string) (*http.Response, error) {
	b, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(c, method, url, bytes.NewBuffer(b))

	if err != nil {
		return nil, err
	}

	for k, v := range header {
		request.Header.Set(k, v)
	}

	client := &http.Client{}
	return client.Do(request)
}

// StringToCapital capitalizes the first character of a string
// Returns empty string if input is empty
func StringToCapital(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	return ""
}

// StructToMap converts any struct to a map[string]any using JSON marshaling/unmarshaling
// This preserves JSON tags and handles nested structures
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

// IsStringInSlice checks if a string exists in a slice of strings
// Returns false if the input string is empty
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
