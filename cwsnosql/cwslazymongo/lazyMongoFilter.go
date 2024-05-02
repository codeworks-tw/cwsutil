/*
 * File: LazyMongoFilter.go
 * Created Date: Wednesday, May 1st 2024, 8:18:43 am
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

type LazyMongoFilter primitive.M

func (f LazyMongoFilter) Eq(key string, value any) LazyMongoFilter {
	f[key] = primitive.M{"$eq": value}
	return f
}

func (f LazyMongoFilter) Ne(key string, value any) LazyMongoFilter {
	f[key] = primitive.M{"$ne": value}
	return f
}

func (f LazyMongoFilter) Gt(key string, value any) LazyMongoFilter {
	f[key] = primitive.M{"$gt": value}
	return f
}

func (f LazyMongoFilter) Gte(key string, value any) LazyMongoFilter {
	f[key] = primitive.M{"$gte": value}
	return f
}

func (f LazyMongoFilter) Lt(key string, value any) LazyMongoFilter {
	f[key] = primitive.M{"$lt": value}
	return f
}

func (f LazyMongoFilter) Lte(key string, value any) LazyMongoFilter {
	f[key] = primitive.M{"$lte": value}
	return f
}

func (f LazyMongoFilter) In(key string, values ...any) LazyMongoFilter {
	a := make(primitive.A, len(values))
	copy(a, values)
	f[key] = primitive.M{"$in": a}
	return f
}

func (f LazyMongoFilter) Nin(key string, values ...any) LazyMongoFilter {
	a := make(primitive.A, len(values))
	copy(a, values)
	f[key] = primitive.M{"$nin": a}
	return f
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
	temp := primitive.A{}
	for _, v := range filters {
		temp = append(temp, v)
	}
	return LazyMongoFilter{"$and": temp}
}

func Or(filters ...LazyMongoFilter) LazyMongoFilter {
	temp := primitive.A{}
	for _, v := range filters {
		temp = append(temp, v)
	}
	return LazyMongoFilter{"$or": temp}
}

func Nor(filters ...LazyMongoFilter) LazyMongoFilter {
	temp := primitive.A{}
	for _, v := range filters {
		temp = append(temp, v)
	}
	return LazyMongoFilter{"$nor": temp}
}

func Not(filter LazyMongoFilter) LazyMongoFilter {
	r := LazyMongoFilter{}
	for k, v := range filter {
		r[k] = primitive.M{"$not": v}
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
		d := make(primitive.M, len(element.(LazyMongoFilter)))
		for k, v := range element.(LazyMongoFilter) {
			d[k] = buildHelper(v)
		}
		return d
	case LazyMongoUpdater:
		d := make(primitive.M, len(element.(LazyMongoUpdater)))
		for k, v := range element.(LazyMongoUpdater) {
			d[k] = buildHelper(v)
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
