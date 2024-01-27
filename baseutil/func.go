/*
 * File: func.go
 * Created Date: Friday, January 26th 2024, 9:49:36 am
 *
 * Last Modified: Sat Jan 27 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package baseutil

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
	"strconv"
	"strings"
	"time"
)

type EnvironmentInfo struct {
	Env       string
	DebugMode bool
	IsLocal   bool
}

func GetEnvironmentInfo() EnvironmentInfo {
	return EnvironmentInfo{
		Env:       GetEnv[string]("ENV"),
		DebugMode: GetEnv("DEBUG", false),
		IsLocal:   GetEnv("IS_LOCAL", false),
	}
}

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

func IntToDateTime(unixTime int64) string {
	return time.Unix(unixTime, 0).Format("2006-01-02 15:04:05")
}

func GetCurrentTimestampString() string {
	return strconv.FormatInt(time.Now().UTC().Unix(), 10)
}

func GetCurrentTimestamp() int64 {
	return time.Now().UTC().Unix()
}

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

func MaxInt32(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func AvoidNilString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ConvertArrayAnyToArrayString(arr []any) []string {
	var result []string
	for _, v := range arr {
		result = append(result, v.(string))
	}
	return result
}

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

func StringToCapital(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	return ""
}
