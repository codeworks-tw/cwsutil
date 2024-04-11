/*
 * File: dynamodb_test.go
 * Created Date: Friday, January 26th 2024, 9:49:36 am
 *
 * Last Modified: Thu Apr 11 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsaws

import (
	"context"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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

func TestDynamoDBProxy(t *testing.T) {
	credential := config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("nve5r", "1wpwhr", ""))

	endPoint := config.WithEndpointResolverWithOptions(
		aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           "http://localhost:8000/",
				SigningRegion: "localhost",
			}, nil
		}))

	TestTableProxy := GetDynamoDBTableProxy[TestItem]("test", context.TODO(), credential, endPoint)

	if TestTableProxy.ProxyTableIsActive() {
		_, err := TestTableProxy.ProxyDeleteTableAndWait()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	if !TestTableProxy.ProxyTableIsActive() {
		_, err := TestTableProxy.ProxyCreateTableAndWaitActive(&dynamodb.CreateTableInput{
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
			log.Fatal(err.Error())
		}
	}

	item, err := TestTableProxy.ProxyGetItem(&dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"Id":    &types.AttributeValueMemberS{Value: "abc"},
			"Item1": &types.AttributeValueMemberN{Value: "123"},
		},
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	if item == nil {
		m := make(map[string]TestItemSub1)
		m["sub1"] = TestItemSub1{
			K1: "This is K1",
			K2: 555,
			K3: 123.456,
		}

		l := []TestItemSub1{}
		l = append(l, TestItemSub1{
			K1: "This is List K1",
			K2: 333,
			K3: 222.777,
		})

		item := TestItem{
			Id:     "abc",
			Item1:  123,
			Item2:  "xxx",
			Array1: []string{"a", "b", "c"},
			List1:  l,
			Map1:   m,
		}

		item2 := TestItem{
			Id:     "def",
			Item1:  123,
			Item2:  "xxx",
			Array1: []string{"d", "e", "f"},
			List1:  l,
			Map1:   m,
		}

		item3 := TestItem{
			Id:     "hij",
			Item1:  123,
			Item2:  "xxx",
			Array1: []string{"h", "i", "j"},
			List1:  l,
			Map1:   m,
		}

		data, err := attributevalue.MarshalMap(item)
		if err != nil {
			log.Fatal(err.Error())
		}

		rc, err := TestTableProxy.ProxyPutItem(&dynamodb.PutItemInput{
			Item: data,
		})

		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println("Put Item: ", rc)

		data, err = attributevalue.MarshalMap(item2)
		if err != nil {
			log.Fatal(err.Error())
		}

		rc, err = TestTableProxy.ProxyPutItem(&dynamodb.PutItemInput{
			Item: data,
		})

		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println("Put Item: ", rc)

		data, err = attributevalue.MarshalMap(item3)
		if err != nil {
			log.Fatal(err.Error())
		}

		rc, err = TestTableProxy.ProxyPutItem(&dynamodb.PutItemInput{
			Item: data,
		})

		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println("Put Item: ", rc)
	}

	item, err = TestTableProxy.ProxyGetItem(&dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"Id":    &types.AttributeValueMemberS{Value: "abc"},
			"Item1": &types.AttributeValueMemberN{Value: "123"},
		},
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Get Item: ", item)

	item, err = TestTableProxy.ProxyUpdateItem(&dynamodb.UpdateItemInput{
		UpdateExpression: aws.String("set Item2 = :r"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":r": &types.AttributeValueMemberS{Value: "ooo"},
		},
		Key: map[string]types.AttributeValue{
			"Id":    &types.AttributeValueMemberS{Value: "abc"},
			"Item1": &types.AttributeValueMemberN{Value: "123"},
		},
		ReturnValues: types.ReturnValueAllNew,
	})

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Update Item: ", item)

	item, err = TestTableProxy.ProxyDeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"Id":    &types.AttributeValueMemberS{Value: "def"},
			"Item1": &types.AttributeValueMemberN{Value: "123"},
		},
		ReturnValues: types.ReturnValueAllOld,
	})

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Delete Item: ", item)

	items, err := TestTableProxy.ProxyScan(&dynamodb.ScanInput{})
	if err != nil {
		log.Fatal(err.Error())
	}
	for item := range items {
		log.Println("Scan Items: ", item)
	}

	items, err = TestTableProxy.ProxyScan(&dynamodb.ScanInput{
		FilterExpression: aws.String("Item2 = :v1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: "xxx"},
		},
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	for item := range items {
		log.Println("Scan Items: ", item)
	}

	items, err = TestTableProxy.ProxyQuery(&dynamodb.QueryInput{
		KeyConditionExpression: aws.String("Id = :v1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: "abc"},
		},
		ProjectionExpression: aws.String("Id, Item1"),
	})

	if err != nil {
		log.Fatal(err.Error())
	}
	for item := range items {
		log.Println("Scan Items: ", item)
	}

	if TestTableProxy.ProxyTableIsActive() {
		_, err := TestTableProxy.ProxyDeleteTableAndWait()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}
