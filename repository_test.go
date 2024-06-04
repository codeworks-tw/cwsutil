/*
 * File: repository_test.go
 * Created Date: Thursday, April 11th 2024, 10:31:37 am
 *
 * Last Modified: Tue Jun 04 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsutil

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TestItemPKey struct {
	Id    string
	Item1 int32
}

type TestItem struct {
	Id     string
	Item1  int32
	Item2  string
	Array1 []string
	List1  []TestItemSub1
	Map1   map[string]TestItemSub1
}

type TestItemSub1 struct {
	K1 string
	K2 int32
	K3 float32
}

var repo Repository[TestItemPKey] = Repository[TestItemPKey]{
	TableName: "TestTable",
}

func createTestTable(ctx context.Context) error {
	proxy := repo.GetDynamoDBTableProxy(ctx)
	if proxy.ProxyTableIsActive() {
		_, err := proxy.ProxyDeleteTableAndWait()
		if err != nil {
			return err
		}
	}

	if !proxy.ProxyTableIsActive() {
		_, err := proxy.ProxyCreateTableAndWaitActive(&dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("Id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("Item1"),
					AttributeType: types.ScalarAttributeTypeN,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("Id"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("Item1"),
					KeyType:       types.KeyTypeRange,
				},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(3),
				WriteCapacityUnits: aws.Int64(3),
			},
			BillingMode: types.BillingModeProvisioned,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func createGSI(ctx context.Context, name string, key string, projections []string) error {
	proxy := repo.GetDynamoDBTableProxy(ctx)

	updateTable := dynamodb.UpdateTableInput{
		TableName: aws.String(proxy.TableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: &key,
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{
			{
				Create: &types.CreateGlobalSecondaryIndexAction{
					IndexName: &name,
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: &key,
							KeyType:       types.KeyTypeHash,
						},
					},
					Projection: &types.Projection{
						ProjectionType:   types.ProjectionTypeInclude,
						NonKeyAttributes: projections,
					},
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(3),
						WriteCapacityUnits: aws.Int64(2),
					},
				},
			},
		},
	}

	_, err := proxy.UpdateTable(*proxy.Context, &updateTable)
	if err != nil {
		return err
	}
	return nil
}

func deleteTestTable(ctx context.Context) {
	proxy := repo.GetDynamoDBTableProxy(ctx)
	if proxy.ProxyTableIsActive() {
		_, err := proxy.ProxyDeleteTableAndWait()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}

func TestRepository(t *testing.T) {
	fmt.Println("\n================ Testing repository ================")

	os.Setenv("ENV", "test")
	os.Setenv("DEBUG", "true")
	os.Setenv("IS_LOCAL", "true")
	os.Setenv("Local_DynamoDB_AWS_ID", "nve5r")
	os.Setenv("Local_DynamoDB_AWS_Secret", "1wpwhr")
	os.Setenv("Local_DynamoDB_URL", "http://localhost:8000/")
	os.Setenv("Local_DynamoDB_REGION", "localhost")

	ctx := context.TODO()
	err := createTestTable(ctx)
	if err != nil {
		log.Println("Local DynamoDB does not exist. Test is skiped.")
		return
	}

	update := expression.Set(expression.Name("Item2"), expression.Value("def"))
	update = update.Set(expression.Name("Array1"), expression.Value([]string{"a", "b", "c"}))
	update = update.Set(expression.Name("List1"), expression.Value([]TestItemSub1{
		{
			K1: "a",
			K2: 1,
			K3: 1.1,
		},
	}))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		t.Error(err)
	}

	item, err := repo.Merge(ctx, TestItemPKey{Id: "abc", Item1: 1}, expr)
	if err != nil {
		log.Fatalln(err)
	}

	if item == nil {
		t.Error("item is nil")
	}

	item, err = repo.Get(ctx, TestItemPKey{Id: "abc", Item1: 1})
	if err != nil {
		log.Fatalln(err)
	}

	if item == nil {
		t.Error("item is nil")
	}

	// create gsi
	indexName := "gsi1"
	err = createGSI(ctx, indexName, "Item2", []string{"Array1", "List1"})
	if err != nil {
		log.Fatalln(err)
	}

	keyexpr := expression.Key("Item2").Equal(expression.Value("def"))
	expr, err = expression.NewBuilder().WithKeyCondition(keyexpr).Build()
	if err != nil {
		t.Error(err)
	}

	items, err := repo.Query(ctx, indexName, expr)
	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range items {
		t.Log(item)
	}

	deleteTestTable(ctx)

	fmt.Println("\n================ Testing repository end ================")
}
