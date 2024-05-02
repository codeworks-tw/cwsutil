/*
 * File: repoEs.go
 * Created Date: Tuesday, April 30th 2024, 8:17:31 pm
 *
 * Last Modified: Thu May 02 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwslazymongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LazyMongoUpdater primitive.M

func (update LazyMongoUpdater) Set(doc any) LazyMongoUpdater {
	update["$set"] = doc
	return update
}

func (update LazyMongoUpdater) AddToSet(key string, values ...any) LazyMongoUpdater {
	update["$addToSet"] = primitive.M{key: primitive.M{"$each": values}}
	return update
}

func (update LazyMongoUpdater) Pull(key string, values ...any) LazyMongoUpdater {
	update["$pull"] = primitive.M{key: primitive.M{"$in": values}}
	return update
}

func (update LazyMongoUpdater) Push(key string, values ...any) LazyMongoUpdater {
	update["$push"] = primitive.M{key: primitive.M{"$each": values}}
	return update
}

func (update LazyMongoUpdater) Pop(key string, head bool) LazyMongoUpdater {
	v := 1
	if head {
		v = -1
	}
	update["$pop"] = primitive.M{key: v}
	return update
}

func (uE LazyMongoUpdater) Build() any {
	return buildHelper(uE)
}

func Set(doc any) LazyMongoUpdater {
	return LazyMongoUpdater{}.Set(doc)
}

func AddToSet(key string, values ...any) LazyMongoUpdater {
	return LazyMongoUpdater{}.AddToSet(key, values...)
}

func Pull(key string, values ...any) LazyMongoUpdater {
	return LazyMongoUpdater{}.Pull(key, values...)
}

func Push(key string, values ...any) LazyMongoUpdater {
	return LazyMongoUpdater{}.Push(key, values...)
}

func Pop(key string, head bool) LazyMongoUpdater {
	return LazyMongoUpdater{}.Pop(key, head)
}
