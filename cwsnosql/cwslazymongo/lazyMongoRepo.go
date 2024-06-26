/*
 * File: lazyMongoRepo.go
 * Created Date: Sunday, May 5th 2024, 5:13:35 pm
 *
 * Last Modified: Tue Jun 04 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwslazymongo

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *LazyMongoRepository) CreateSimpleAscendingIndex(ctx context.Context, key string, opts ...*options.CreateIndexesOptions) (string, error) {
	model := mongo.IndexModel{
		Keys: LazyMongoFilter{primitive.E{Key: key, Value: 1}}.Build(),
	}
	return r.CreateIndex(ctx, model, opts...)
}

func (r *LazyMongoRepository) CreateSimpleUniqueAscendingIndex(ctx context.Context, key string, opts ...*options.CreateIndexesOptions) (string, error) {
	model := mongo.IndexModel{
		Keys:    LazyMongoFilter{primitive.E{Key: key, Value: 1}}.Build(),
		Options: options.Index().SetUnique(true),
	}
	return r.CreateIndex(ctx, model, opts...)
}

func (r *LazyMongoRepository) CreateIndex(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return "", err
	}
	return collection.Indexes().CreateOne(ctx, model, opts...)
}

func (r *LazyMongoRepository) DeleteIndex(ctx context.Context, indexName string, opts ...*options.DropIndexesOptions) (bson.Raw, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.Indexes().DropOne(ctx, indexName, opts...)
}

func (r *LazyMongoRepository) Get(ctx context.Context, filter LazyMongoFilter, out any, opts ...*options.FindOneOptions) error {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return err
	}
	return collection.FindOne(ctx, filter.Build(), opts...).Decode(out)
}

func (r *LazyMongoRepository) Select(ctx context.Context, filter LazyMongoFilter, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.Find(ctx, filter.Build(), opts...)
}

func (r *LazyMongoRepository) Add(ctx context.Context, data any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.InsertOne(ctx, data, opts...)
}

func (r *LazyMongoRepository) AddMany(ctx context.Context, data []any, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.InsertMany(ctx, data, opts...)
}

func (r *LazyMongoRepository) Update(ctx context.Context, filter LazyMongoFilter, update LazyMongoUpdater, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.UpdateOne(ctx, filter.Build(), update.Build(), opts...)
}

func (r *LazyMongoRepository) UpdateMany(ctx context.Context, filter LazyMongoFilter, update LazyMongoUpdater, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.UpdateMany(ctx, filter.Build(), update.Build(), opts...)
}

func (r *LazyMongoRepository) Delete(ctx context.Context, filter LazyMongoFilter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.DeleteOne(ctx, filter.Build(), opts...)
}

func (r *LazyMongoRepository) DeleteMany(ctx context.Context, filter LazyMongoFilter, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	return collection.DeleteMany(ctx, filter.Build(), opts...)
}

func (r *LazyMongoRepository) Exist(ctx context.Context, filter LazyMongoFilter, opts ...*options.CountOptions) (bool, error) {
	count, err := r.Count(ctx, filter, opts...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *LazyMongoRepository) Count(ctx context.Context, filter LazyMongoFilter, opts ...*options.CountOptions) (int64, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return 0, err
	}
	return collection.CountDocuments(ctx, filter.Build(), opts...)
}

func (r *LazyMongoRepository) Upsert(ctx context.Context, filter LazyMongoFilter, update LazyMongoUpdater) (*mongo.UpdateResult, error) {
	return r.Update(ctx, filter, update, options.Update().SetUpsert(true))
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
