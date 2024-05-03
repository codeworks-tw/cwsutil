/*
 * File: LazyMongoFilter.go
 * Created Date: Wednesday, May 1st 2024, 8:18:43 am
 *
 * Last Modified: Fri May 03 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwslazymongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LazyMongoFilter primitive.D

func (f LazyMongoFilter) Eq(key string, value any) LazyMongoFilter {
	return append(f, primitive.E{Key: key, Value: primitive.D{primitive.E{Key: "$eq", Value: value}}})
}

func (f LazyMongoFilter) Ne(key string, value any) LazyMongoFilter {
	return append(f, primitive.E{Key: key, Value: primitive.D{primitive.E{Key: "$ne", Value: value}}})
}

func (f LazyMongoFilter) Gt(key string, value any) LazyMongoFilter {
	return append(f, primitive.E{Key: key, Value: primitive.D{primitive.E{Key: "$gt", Value: value}}})
}

func (f LazyMongoFilter) Gte(key string, value any) LazyMongoFilter {
	return append(f, primitive.E{Key: key, Value: primitive.D{primitive.E{Key: "$gte", Value: value}}})
}

func (f LazyMongoFilter) Lt(key string, value any) LazyMongoFilter {
	return append(f, primitive.E{Key: key, Value: primitive.D{primitive.E{Key: "$lt", Value: value}}})
}

func (f LazyMongoFilter) Lte(key string, value any) LazyMongoFilter {
	return append(f, primitive.E{Key: key, Value: primitive.D{primitive.E{Key: "$lte", Value: value}}})
}

func (f LazyMongoFilter) In(key string, values ...any) LazyMongoFilter {
	return append(f, primitive.E{Key: key, Value: primitive.D{primitive.E{Key: "$in", Value: values}}})
}

func (f LazyMongoFilter) Nin(key string, values ...any) LazyMongoFilter {
	return append(f, primitive.E{Key: key, Value: primitive.D{primitive.E{Key: "$nin", Value: values}}})
}

func All() LazyMongoFilter {
	return LazyMongoFilter{}
}

func Eq(key string, value any) LazyMongoFilter {
	return LazyMongoFilter{}.Eq(key, value)
}

func Ne(key string, value any) LazyMongoFilter {
	return LazyMongoFilter{}.Ne(key, value)
}

func Gt(key string, value any) LazyMongoFilter {
	return LazyMongoFilter{}.Gt(key, value)
}

func Gte(key string, value any) LazyMongoFilter {
	return LazyMongoFilter{}.Gte(key, value)
}

func Lt(key string, value any) LazyMongoFilter {
	return LazyMongoFilter{}.Lt(key, value)
}

func Lte(key string, value any) LazyMongoFilter {
	return LazyMongoFilter{}.Lte(key, value)
}

func In(key string, values ...any) LazyMongoFilter {
	return LazyMongoFilter{}.In(key, values...)
}

func Nin(key string, values ...any) LazyMongoFilter {
	return LazyMongoFilter{}.Nin(key, values...)
}

func And(filters ...LazyMongoFilter) LazyMongoFilter {
	return LazyMongoFilter{
		primitive.E{Key: "$and", Value: filters},
	}
}

func Or(filters ...LazyMongoFilter) LazyMongoFilter {
	return LazyMongoFilter{
		primitive.E{Key: "$or", Value: filters},
	}
}

func Nor(filters ...LazyMongoFilter) LazyMongoFilter {
	return LazyMongoFilter{
		primitive.E{Key: "$nor", Value: filters},
	}
}

func Not(filter LazyMongoFilter) LazyMongoFilter {
	r := make(LazyMongoFilter, len(filter))
	for i, e := range filter {
		r[i] = primitive.E{Key: e.Key, Value: primitive.D{primitive.E{Key: "$not", Value: e.Value}}}
	}
	return r
}

func (uE LazyMongoFilter) Build() any {
	return buildHelper(uE)
}

func buildHelper(element any) any {
	switch element.(type) {
	case primitive.A:
		for i, e := range element.(primitive.A) {
			element.(primitive.A)[i] = buildHelper(e)
		}
	case LazyMongoFilter:
		d := primitive.D{}
		for _, v := range element.(LazyMongoFilter) {
			d = append(d, buildHelper(v).(primitive.E))
		}
		return d
	case LazyMongoUpdater:
		d := primitive.D{}
		for _, v := range element.(LazyMongoUpdater) {
			d = append(d, buildHelper(v).(primitive.E))
		}
		return d
	case primitive.D:
		for i, e := range element.(primitive.D) {
			element.(primitive.D)[i] = buildHelper(e).(primitive.E)
		}
	case primitive.M:
		for k, v := range element.(primitive.M) {
			element.(primitive.M)[k] = buildHelper(v)
		}
	case primitive.E:
		e := element.(primitive.E)
		e.Value = buildHelper(e.Value)
		return e
	}
	return element
}
