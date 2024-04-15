/*
 * File: mongo.go
 * Created Date: Thursday, April 11th 2024, 3:11:23 pm
 *
 * Last Modified: Mon Apr 15 2024
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

	for k := range keys {
		keys[k] = 1 // set to ascending
	}

	c, err := collection.Indexes().List(ctx)
	if err != nil {
		return err
	}
	for c.TryNext(ctx) {
		var m map[string]any
		index := c.Current
		index.Lookup("key").Unmarshal(&m)
		check := true
		for k := range m {
			if _, ok := keys[k]; !ok {
				check = false
				break
			}
		}
		if check {
			return nil
		}
	}

	// create index
	model := mongo.IndexModel{
		Keys:    keys,
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
