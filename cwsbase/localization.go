/*
 * File: localization.go
 * Created Date: Tuesday, April 23rd 2024, 2:51:04 pm
 *
 * Last Modified: Tue Jun 04 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsbase

import (
	"encoding/json"
	"fmt"
	"sync"
)

type LocalizationCode string
type LocalizationLanguage string

const (
	English   LocalizationLanguage = "en"
	Taiwanese LocalizationLanguage = "zh_tw"
	Chinese   LocalizationLanguage = "zh_cn"
)

var localInitLock sync.Mutex
var localmap map[LocalizationLanguage]map[LocalizationCode]string = map[LocalizationLanguage]map[LocalizationCode]string{}

const (
	LocalCode_InternalServerError LocalizationCode = "500"
	LocalCode_Unauthorized        LocalizationCode = "401"
	LocalCode_BadRequest          LocalizationCode = "400"
	LocalCode_OK                  LocalizationCode = "200"
	LocalCode_Forbidden           LocalizationCode = "403"
	LocalCode_NotFound            LocalizationCode = "404"
)

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

func InitBasicLocalizationData() {
	UpdateLocalizationData([]byte(localdata))
}

func UpdateLocalizationData(jsonData []byte) error {
	localInitLock.Lock()
	defer localInitLock.Unlock()

	var data map[LocalizationLanguage]map[LocalizationCode]string
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}

	for k, v := range data {
		if val, ok := localmap[k]; ok {
			for k2, v2 := range v {
				val[k2] = v2
			}
		} else {
			localmap[k] = v
		}
	}
	return nil
}

func GetLocalizationMessage(code LocalizationCode, strs ...any) string {
	lang := GetEnv("LOCALIZATION_LANGUAGE", "en")

	if langMap, ok := localmap[LocalizationLanguage(lang)]; ok {
		if message, ok := langMap[LocalizationCode(code)]; ok {
			if len(strs) == 0 {
				return message
			}
			return fmt.Sprintf(message, strs...)
		}
	}

	return fmt.Sprintf("Localization code %s does not exist in language %s ", code, lang)
}
