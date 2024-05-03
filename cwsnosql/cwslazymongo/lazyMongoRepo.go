/*
 * File: repo.go
 * Created Date: Tuesday, April 30th 2024, 4:44:31 pm
 *
 * Last Modified: Fri May 03 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwslazymongo

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var lock sync.Mutex = sync.Mutex{}

var clients map[string]*mongo.Client = map[string]*mongo.Client{}

func GetMongoSingletonClient(url string, ctx context.Context) (*mongo.Client, error) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := clients[url]; !ok {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
		if err != nil {
			return client, err
		}
		clients[url] = client
	}
	return clients[url], nil
}

func CloseSingletonClient(url string, ctx context.Context) {
	lock.Lock()
	defer lock.Unlock()

	if client, ok := clients[url]; ok {
		client.Disconnect(ctx)
		delete(clients, url)
	}
}

type LazyMongoRepository struct {
	Url            string
	DbName         string
	CollectionName string
}

func (r *LazyMongoRepository) GetCollection(ctx context.Context) (*mongo.Collection, error) {
	client, err := GetMongoSingletonClient(r.Url, ctx)
	if err != nil {
		return nil, err
	}
	return client.Database(r.DbName).Collection(r.CollectionName), nil
}

func (r *LazyMongoRepository) Get(ctx context.Context, filter LazyMongoFilter, out any) error {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return err
	}
	return collection.FindOne(ctx, filter.Build()).Decode(out)
}

func (r *LazyMongoRepository) Select(ctx context.Context, filter LazyMongoFilter) (*mongo.Cursor, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.Find(ctx, filter.Build())
}

func (r *LazyMongoRepository) Add(ctx context.Context, data any) (*mongo.InsertOneResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.InsertOne(ctx, data)
}

func (r *LazyMongoRepository) AddMany(ctx context.Context, data []any) (*mongo.InsertManyResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.InsertMany(ctx, data)
}

func (r *LazyMongoRepository) Update(ctx context.Context, filter LazyMongoFilter, update LazyMongoUpdater) (*mongo.UpdateResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.UpdateOne(ctx, filter.Build(), update.Build())
}

func (r *LazyMongoRepository) UpdateMany(ctx context.Context, filter LazyMongoFilter, update LazyMongoUpdater) (*mongo.UpdateResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.UpdateMany(ctx, filter.Build(), update.Build())
}

func (r *LazyMongoRepository) Delete(ctx context.Context, filter LazyMongoFilter) (*mongo.DeleteResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.DeleteOne(ctx, filter.Build())
}

func (r *LazyMongoRepository) DeleteMany(ctx context.Context, filter LazyMongoFilter) (*mongo.DeleteResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.DeleteMany(ctx, filter.Build())
}

func (r *LazyMongoRepository) Exist(ctx context.Context, filter LazyMongoFilter) (bool, error) {
	count, err := r.Count(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *LazyMongoRepository) Count(ctx context.Context, filter LazyMongoFilter) (int64, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return 0, err
	}
	return collection.CountDocuments(ctx, filter.Build())
}

func (r *LazyMongoRepository) Upsert(ctx context.Context, filter LazyMongoFilter, update LazyMongoUpdater) (*mongo.UpdateResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.UpdateOne(ctx, filter.Build(), update.Build(), options.Update().SetUpsert(true))
}

func MarshalToFilter(data any) (LazyMongoFilter, error) {
	m := LazyMongoFilter{}
	dataBytes, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(dataBytes, &m)
	return m, err
}

func MarshalToUpdater(data any) LazyMongoUpdater {
	return Set(data)
}
