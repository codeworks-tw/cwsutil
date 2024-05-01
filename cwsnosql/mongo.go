/*
 * File: mongo.go
 * Created Date: Thursday, April 11th 2024, 3:11:23 pm
 *
 * Last Modified: Mon Apr 29 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsnosql

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

type MongoDBRepository[PKey any] struct {
	Url            string
	DbName         string
	CollectionName string
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

func (r *MongoDBRepository[PKey]) GetCollection(ctx context.Context) (*mongo.Collection, error) {
	client, err := GetMongoSingletonClient(r.Url, ctx)
	if err != nil {
		return nil, err
	}
	return client.Database(r.DbName).Collection(r.CollectionName), nil
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
		temp = append(temp, bson.E{Key: k, Value: keys[k]})
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
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return err
	}

	filter, err := marshalToBsonMap(pkey)
	if err != nil {
		return err
	}

	update := bson.D{{"$set", doc}}
	_, err = collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *MongoDBRepository[PKey]) AddValuesToSet(ctx context.Context, pkey PKey, key string, values ...any) (*mongo.UpdateResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	filter, err := marshalToBsonMap(pkey)
	if err != nil {
		return nil, err
	}

	update := bson.D{{"$addToSet", map[string]any{key: map[string][]any{"$each": values}}}}
	return collection.UpdateOne(ctx, filter, update)
}

func (r *MongoDBRepository[PKey]) PullValuesFromSet(ctx context.Context, pkey PKey, key string, values ...any) (*mongo.UpdateResult, error) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	filter, err := marshalToBsonMap(pkey)
	if err != nil {
		return nil, err
	}

	update := bson.D{{"$pull", map[string]any{key: map[string][]any{"$in": values}}}}
	return collection.UpdateOne(ctx, filter, update)
}

func (r *MongoDBRepository[PKey]) Find(ctx context.Context, pkey PKey, out any) error {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return err
	}

	filter, err := marshalToBsonMap(pkey)
	if err != nil {
		return err
	}

	err = collection.FindOne(ctx, filter).Decode(out)
	return err
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
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return err
	}

	filter, err := marshalToBsonMap(pkey)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}