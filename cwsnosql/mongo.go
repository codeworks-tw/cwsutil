/*
 * File: mongo.go
 * Created Date: Thursday, April 11th 2024, 3:11:23 pm
 *
 * Last Modified: Thu May 02 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsnosql

import (
	"context"
	"sync"

	"github.com/codeworks-tw/cwsutil/cwsnosql/cwslazymongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRepository[PKey any] struct {
	once           sync.Once
	lazyRepo       cwslazymongo.LazyMongoRepository
	Url            string
	DbName         string
	CollectionName string
}

func (r *MongoDBRepository[PKey]) ToLazyMongoRepository() *cwslazymongo.LazyMongoRepository {
	r.once.Do(func() {
		r.lazyRepo = cwslazymongo.LazyMongoRepository{
			Url:            r.Url,
			DbName:         r.DbName,
			CollectionName: r.CollectionName,
		}
	})
	return &r.lazyRepo
}

func (r *MongoDBRepository[PKey]) GetCollection(ctx context.Context) (*mongo.Collection, error) {
	return r.ToLazyMongoRepository().GetCollection(ctx)
}

func marshalToBsonMap(data any) (bson.M, error) {
	m := bson.M{}
	dataBytes, err := bson.Marshal(data)
	if err != nil {
		return m, err
	}
	err = bson.Unmarshal(dataBytes, &m)
	return m, err
}

func (r *MongoDBRepository[PKey]) CreateSimpleUniqueAscendingIndex(ctx context.Context) error {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return err
	}

	var pkey PKey
	keys, err := marshalToBsonMap(pkey)
	if err != nil {
		return err
	}

	temp := bson.D{}
	for k := range keys {
		temp = append(temp, bson.E{Key: k, Value: 1})
	}

	// create index
	model := mongo.IndexModel{
		Keys:    temp,
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(ctx, model)
	return err
}

func (r *MongoDBRepository[PKey]) Upsert(ctx context.Context, pkey PKey, doc any) error {
	filter, err := cwslazymongo.MarshalToFilter(pkey)
	if err != nil {
		return err
	}

	_, err = r.ToLazyMongoRepository().Upsert(ctx, filter, cwslazymongo.Set((doc)))
	return err
}

func (r *MongoDBRepository[PKey]) AddValuesToSet(ctx context.Context, pkey PKey, key string, values ...any) (*mongo.UpdateResult, error) {
	filter, err := cwslazymongo.MarshalToFilter(pkey)
	if err != nil {
		return nil, err
	}
	return r.ToLazyMongoRepository().Update(ctx, filter, cwslazymongo.AddToSet(key, values...))
}

func (r *MongoDBRepository[PKey]) PullValuesFromSet(ctx context.Context, pkey PKey, key string, values ...any) (*mongo.UpdateResult, error) {
	filter, err := cwslazymongo.MarshalToFilter(pkey)
	if err != nil {
		return nil, err
	}
	return r.ToLazyMongoRepository().Update(ctx, filter, cwslazymongo.Pull(key, values...))
}

func (r *MongoDBRepository[PKey]) Find(ctx context.Context, pkey PKey, out any) error {
	filter, err := cwslazymongo.MarshalToFilter(pkey)
	if err != nil {
		return err
	}
	return r.ToLazyMongoRepository().Get(ctx, filter, out)
}

func (r *MongoDBRepository[PKey]) FindWithFilter(ctx context.Context, filter bson.M, out any) error {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return err
	}

	err = collection.FindOne(ctx, filter).Decode(out)
	return err
}

func (r *MongoDBRepository[PKey]) Exist(ctx context.Context, filter bson.M) (bool, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return false, err
	}

	//filter, err := marshalToBsonMap(pkey)
	//if err != nil {
	//	return false, err
	//}

	count, err := collection.CountDocuments(ctx, filter, nil)

	if count > 0 {
		return true, nil
	}
	return false, err
}

func (r *MongoDBRepository[PKey]) Delete(ctx context.Context, pkey PKey) error {
	filter, err := cwslazymongo.MarshalToFilter(pkey)
	if err != nil {
		return err
	}
	_, err = r.ToLazyMongoRepository().Delete(ctx, filter)
	return err
}
