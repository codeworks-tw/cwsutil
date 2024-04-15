/*
 * File: nosql_test.go
 * Created Date: Friday, April 12th 2024, 4:45:03 pm
 *
 * Last Modified: Mon Apr 15 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsutil

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/codeworks-tw/cwsutil/cwsnosql"
)

type NoSqlTestItemKey struct {
	Id string `bson:"id"`
}

type NoSqlTestItem struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}

var RepositoryNoSqlTest = &cwsnosql.MongoDBRepository[NoSqlTestItemKey]{
	Url:            "mongodb://localhost:27017",
	DbName:         "test",
	CollectionName: "test",
}

func TestMongo(t *testing.T) {
	fmt.Println("\n================ Testing nosql repository ================")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pkey := NoSqlTestItemKey{
		Id: "1",
	}

	item := NoSqlTestItem{
		Id:   "1",
		Name: "testname",
	}

	err := RepositoryNoSqlTest.Upsert(ctx, pkey, item)
	if err != nil {
		log.Println(err)
		return
	}

	var data NoSqlTestItem
	err = RepositoryNoSqlTest.Find(ctx, pkey, &data)
	if err != nil {
		log.Println(err)
		return
	}

	err = RepositoryNoSqlTest.Delete(ctx, pkey)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("\n================ Testing nosql repository end ================")
}
