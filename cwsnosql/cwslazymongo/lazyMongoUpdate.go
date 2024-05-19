/*
 * File: repoEs.go
 * Created Date: Tuesday, April 30th 2024, 8:17:31 pm
 *
 * Last Modified: Sun May 19 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwslazymongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LazyMongoUpdater primitive.D

func (update LazyMongoUpdater) SetKey(key string, value any) LazyMongoUpdater {
	return append(update, primitive.E{Key: "$set", Value: primitive.D{primitive.E{Key: key, Value: value}}})
}

func (update LazyMongoUpdater) Inc(key string, value int) LazyMongoUpdater {
	return append(update, primitive.E{Key: "$inc", Value: primitive.D{primitive.E{Key: key, Value: value}}})
}

func (update LazyMongoUpdater) Set(doc any) LazyMongoUpdater {
	return append(update, primitive.E{Key: "$set", Value: doc})
}

func (update LazyMongoUpdater) Push(key string, values ...any) LazyMongoUpdater {
	return append(update, primitive.E{Key: "$push", Value: primitive.D{
		primitive.E{Key: key, Value: primitive.D{primitive.E{Key: "$each", Value: values}}},
	}})
}

func (update LazyMongoUpdater) AddToSet(key string, values ...any) LazyMongoUpdater {
	return append(update, primitive.E{Key: "$addToSet", Value: primitive.D{
		primitive.E{Key: key, Value: primitive.D{primitive.E{Key: "$each", Value: values}}},
	}})
}

func (update LazyMongoUpdater) Pull(key string, values ...any) LazyMongoUpdater {
	return append(update, primitive.E{Key: "$pull", Value: primitive.D{
		primitive.E{Key: key, Value: primitive.D{primitive.E{Key: "$in", Value: values}}},
	}})
}

func (update LazyMongoUpdater) Pop(key string, head bool) LazyMongoUpdater {
	v := 1
	if head {
		v = -1
	}
	return append(update, primitive.E{Key: "$pop", Value: primitive.D{
		primitive.E{Key: key, Value: v},
	}})
}

func (uE LazyMongoUpdater) Build() any {
	return buildHelper(uE)
}

func Set(doc any) LazyMongoUpdater {
	return LazyMongoUpdater{}.Set(doc)
}

func Inc(key string, value int) LazyMongoUpdater {
	return LazyMongoUpdater{}.Inc(key, value)
}

func SetKey(key string, value any) LazyMongoUpdater {
	return LazyMongoUpdater{}.SetKey(key, value)
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
