/*
 * File: nosql_test.go
 * Created Date: Friday, April 12th 2024, 4:45:03 pm
 *
 * Last Modified: Wed May 01 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/codeworks-tw/cwsutil/cwsnosql/cwslazymongo"
)

type NoSqlTestItemKey struct {
	Id string `bson:"id"`
}

type NoSqlTestItem struct {
	Id   string   `bson:"id"`
	Name string   `bson:"name"`
	Tags []string `bson:"tags"`
}

var RepositoryNoSqlTest = cwslazymongo.LazyMongoRepository{
	Url:            "mongodb://localhost:27017",
	DbName:         "testnosql",
	CollectionName: "testnosqllazy",
}

func TestMongo(t *testing.T) {
	fmt.Println("\n================ Testing nosql lazy mongo repo ================")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	item := NoSqlTestItem{
		Id:   "1",
		Name: "testname",
		Tags: []string{},
	}

	_, err := RepositoryNoSqlTest.Add(ctx, NoSqlTestItem{
		Id:   "2",
		Name: "testname2",
		Tags: []string{},
	})
	if err != nil {
		t.Error(err)
		return
	}

	_, err = RepositoryNoSqlTest.AddMany(ctx, []any{
		NoSqlTestItem{
			Id:   "3",
			Name: "testname3",
			Tags: []string{},
		},
		NoSqlTestItem{
			Id:   "4",
			Name: "testname4",
			Tags: []string{},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	r, err := RepositoryNoSqlTest.Upsert(ctx, cwslazymongo.Eq("id", "1"), cwslazymongo.Set(item))
	if err != nil {
		t.Error(err)
		return
	}
	if r.UpsertedCount == 0 && r.MatchedCount == 0 {
		t.Error("no match")
		return
	}

	count, err := RepositoryNoSqlTest.Count(ctx, cwslazymongo.Eq("tags", []string{}))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(count)

	if exist, err := RepositoryNoSqlTest.Exist(ctx, cwslazymongo.Eq("id", "1")); err != nil {
		t.Error(err)
		return
	} else if !exist {
		t.Error("no match")
		return
	}

	ur, err := RepositoryNoSqlTest.Update(ctx, cwslazymongo.Eq("id", "1"), cwslazymongo.Set(map[string]any{"name": "new_name"}))
	if err != nil {
		t.Error(err)
		return
	}
	if ur.ModifiedCount == 0 {
		t.Error("no match")
		return
	}

	var data NoSqlTestItem
	err = RepositoryNoSqlTest.Get(ctx, cwslazymongo.Or(cwslazymongo.Eq("id", "1"), cwslazymongo.Eq("name", "testname")), &data)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("Get Item:", data)

	usr, err := RepositoryNoSqlTest.UpdateMany(ctx, cwslazymongo.Eq("tags", []string{}), cwslazymongo.AddToSet("tags", "t1", "t2"))
	if err != nil {
		t.Error(err)
		return
	}
	if usr.MatchedCount == 0 {
		t.Error("no match")
		return
	}

	datas, err := RepositoryNoSqlTest.Select(ctx, cwslazymongo.In("tags", "t1", "t2"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("Select Items: ", len(datas))

	d, err := RepositoryNoSqlTest.Delete(ctx, cwslazymongo.Eq("id", "1"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(d.DeletedCount)

	ds, err := RepositoryNoSqlTest.DeleteMany(ctx, cwslazymongo.Eq("tags", []string{"t1", "t2"}))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(ds.DeletedCount)

	fmt.Println("\n================ Testing nosql repository end ================")
}
