/*
 * File: repoEs.go
 * Created Date: Tuesday, April 30th 2024, 8:17:31 pm
 *
 * Last Modified: Wed May 01 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwslazymongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LazyMongoUpdate primitive.M

func (update LazyMongoUpdate) Set(doc any) LazyMongoUpdate {
	update["$set"] = doc
	return update
}

func (update LazyMongoUpdate) AddToSet(key string, values ...any) LazyMongoUpdate {
	update["$addToSet"] = primitive.M{key: primitive.M{"$each": values}}
	return update
}

func (update LazyMongoUpdate) Pull(key string, values ...any) LazyMongoUpdate {
	update["$pull"] = primitive.M{key: primitive.M{"$in": values}}
	return update
}

func (update LazyMongoUpdate) Push(key string, values ...any) LazyMongoUpdate {
	update["$push"] = primitive.M{key: primitive.M{"$each": values}}
	return update
}

func (update LazyMongoUpdate) Pop(key string, head bool) LazyMongoUpdate {
	v := 1
	if head {
		v = -1
	}
	update["$pop"] = primitive.M{key: v}
	return update
}

func (uE LazyMongoUpdate) Build() any {
	return buildHelper(uE)
}

func Set(doc any) LazyMongoUpdate {
	return LazyMongoUpdate{}.Set(doc)
}

func AddToSet(key string, values ...any) LazyMongoUpdate {
	return LazyMongoUpdate{}.AddToSet(key, values...)
}

func Pull(key string, values ...any) LazyMongoUpdate {
	return LazyMongoUpdate{}.Pull(key, values...)
}

func Push(key string, values ...any) LazyMongoUpdate {
	return LazyMongoUpdate{}.Push(key, values...)
}

func Pop(key string, head bool) LazyMongoUpdate {
	return LazyMongoUpdate{}.Pop(key, head)
}
