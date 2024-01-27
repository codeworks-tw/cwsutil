/*
 * File: dynamodb.go
 * Created Date: Friday, January 26th 2024, 9:49:36 am
 *
 * Last Modified: Fri Jan 26 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 Codeworks Ltd.
 */

package awsutil

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBTableProxy[O any] struct {
	*dynamodb.Client
	Context   *context.Context
	TableName string
}

func GetDynamoDBTableProxy[O any](name string, ctx context.Context, optFns ...func(*config.LoadOptions) error) DynamoDBTableProxy[O] {
	if ctx == nil {
		ctx = context.TODO()
	}

	return DynamoDBTableProxy[O]{
		TableName: name,
		Context:   &ctx,
		Client: GetSingletonClient(ClientName_DynamoDB, ctx, func(cfg aws.Config) *dynamodb.Client {
			return dynamodb.NewFromConfig(cfg)
		}, optFns...),
	}
}

func (table *DynamoDBTableProxy[O]) ProxyTableIsActive() bool {
	out, err := table.DescribeTable(*table.Context, &dynamodb.DescribeTableInput{
		TableName: &table.TableName,
	})
	if err != nil {
		return false
	}
	return out.Table.TableStatus == types.TableStatusActive
}

func (table *DynamoDBTableProxy[O]) ProxyTableHasGSI(gsiName string) bool {
	out, err := table.DescribeTable(*table.Context, &dynamodb.DescribeTableInput{
		TableName: &table.TableName,
	})
	if err != nil {
		return false
	}
	for i := range out.Table.GlobalSecondaryIndexes {
		if *out.Table.GlobalSecondaryIndexes[i].IndexName == gsiName {
			return true
		}
	}
	return false
}

func (table *DynamoDBTableProxy[O]) ProxyCreateTableAndWaitActive(input *dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error) {
	input.TableName = &table.TableName
	out, err := table.Client.CreateTable(*table.Context, input)
	if err != nil {
		return out, err
	}
	w := dynamodb.NewTableExistsWaiter(table.Client)
	err = w.Wait(*table.Context,
		&dynamodb.DescribeTableInput{
			TableName: &table.TableName,
		},
		2*time.Minute,
		func(o *dynamodb.TableExistsWaiterOptions) {
			o.MaxDelay = 5 * time.Second
			o.MinDelay = 5 * time.Second
		})
	if err != nil {
		return out, err
	}
	log.Printf("Table %s has been created and activated", table.TableName)
	return out, err
}

func (table *DynamoDBTableProxy[O]) ProxyDeleteTableAndWait() (*dynamodb.DeleteTableOutput, error) {
	out, err := table.Client.DeleteTable(*table.Context, &dynamodb.DeleteTableInput{
		TableName: &table.TableName,
	})
	if err != nil {
		return out, err
	}
	w := dynamodb.NewTableNotExistsWaiter(table.Client)
	err = w.Wait(*table.Context,
		&dynamodb.DescribeTableInput{
			TableName: &table.TableName,
		},
		2*time.Minute,
		func(o *dynamodb.TableNotExistsWaiterOptions) {
			o.MaxDelay = 5 * time.Second
			o.MinDelay = 5 * time.Second
		})
	if err != nil {
		return out, err
	}
	log.Printf("Table %s has been deleted", table.TableName)
	return out, err
}

func (table *DynamoDBTableProxy[O]) unmarshalMap(item map[string]types.AttributeValue) (*O, error) {
	if item != nil {
		var data O
		err := attributevalue.UnmarshalMap(item, &data)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
	return nil, nil
}

func (table *DynamoDBTableProxy[O]) ProxyPutItem(input *dynamodb.PutItemInput) (*O, error) {
	input.TableName = &table.TableName
	out, err := table.PutItem(*table.Context, input)
	if err != nil {
		return nil, err
	}
	return table.unmarshalMap(out.Attributes)
}

func (table *DynamoDBTableProxy[O]) ProxyGetItem(input *dynamodb.GetItemInput) (*O, error) {
	input.TableName = &table.TableName
	out, err := table.GetItem(*table.Context, input)
	if err != nil {
		return nil, err
	}
	return table.unmarshalMap(out.Item)
}

func (table *DynamoDBTableProxy[O]) ProxyUpdateItem(input *dynamodb.UpdateItemInput) (*O, error) {
	input.TableName = &table.TableName
	out, err := table.UpdateItem(*table.Context, input)
	if err != nil {
		return nil, err
	}
	return table.unmarshalMap(out.Attributes)
}

func (table *DynamoDBTableProxy[O]) ProxyDeleteItem(input *dynamodb.DeleteItemInput) (*O, error) {
	input.TableName = &table.TableName
	out, err := table.DeleteItem(*table.Context, input)
	if err != nil {
		return nil, err
	}
	return table.unmarshalMap(out.Attributes)
}

func (table *DynamoDBTableProxy[O]) ProxyQueryBatchUpdate(input *dynamodb.QueryInput, callback func(item *O, batchRequests *[]types.WriteRequest) *[]types.WriteRequest) *sync.WaitGroup {
	var w sync.WaitGroup
	w.Add(1)
	go func() {
		defer w.Done()
		input.TableName = &table.TableName
		p := dynamodb.NewQueryPaginator(table.Client, input)

		var batchRequests *[]types.WriteRequest = &[]types.WriteRequest{}
		for p.HasMorePages() {
			out, err := p.NextPage(*table.Context)
			if err != nil {
				panic(err.Error())
			} else {
				for _, v := range out.Items {
					var data O
					err = attributevalue.UnmarshalMap(v, &data)
					if err != nil {
						log.Println(err.Error())
					} else if callback != nil {
						batchRequests = callback(&data, batchRequests)
					}
				}
			}
		}
		if batchRequests != nil && len(*batchRequests) > 0 {
			_, e := table.BatchWriteItem(*table.Context, &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]types.WriteRequest{
					table.TableName: *batchRequests,
				},
			})
			if e != nil {
				panic(e.Error())
			}
		}
	}()
	return &w
}

func (table *DynamoDBTableProxy[O]) ProxyScanWithCalback(input *dynamodb.ScanInput, callback func(item *O, batchRequests []types.WriteRequest) []types.WriteRequest) *sync.WaitGroup {
	var w sync.WaitGroup
	w.Add(1)
	go func() {
		defer w.Done()
		input.TableName = &table.TableName
		p := dynamodb.NewScanPaginator(table.Client, input)

		batchRequests := []types.WriteRequest{}
		for p.HasMorePages() {
			out, err := p.NextPage(*table.Context)
			if err != nil {
				panic(err.Error())
			}
			for _, v := range out.Items {
				var data O
				err = attributevalue.UnmarshalMap(v, &data)
				if err != nil {
					log.Println(err.Error())
				} else if callback != nil {
					batchRequests = callback(&data, batchRequests)
				}
			}
		}
		if len(batchRequests) > 0 {
			_, e := table.BatchWriteItem(*table.Context, &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]types.WriteRequest{
					table.TableName: batchRequests,
				},
			})
			if e != nil {
				panic(e.Error())
			}
		}
	}()
	return &w
}

func (table *DynamoDBTableProxy[O]) ProxyQuery(input *dynamodb.QueryInput) ([]*O, error) {
	input.TableName = &table.TableName
	p := dynamodb.NewQueryPaginator(table.Client, input)

	var items []*O = make([]*O, 0)
	for p.HasMorePages() {
		out, err := p.NextPage(*table.Context)
		if err != nil {
			return items, err
		} else {
			for _, v := range out.Items {
				var data O
				err = attributevalue.UnmarshalMap(v, &data)
				if err != nil {
					log.Println(err.Error())
				} else {
					items = append(items, &data)
				}
			}
		}
	}
	return items, nil
}

func (table *DynamoDBTableProxy[O]) ProxyScan(input *dynamodb.ScanInput) ([]*O, error) {
	input.TableName = &table.TableName
	p := dynamodb.NewScanPaginator(table.Client, input)

	var items []*O = make([]*O, 0)
	for p.HasMorePages() {
		out, err := p.NextPage(*table.Context)
		if err != nil {
			return items, err
		}
		for _, v := range out.Items {
			var data O
			err = attributevalue.UnmarshalMap(v, &data)
			if err != nil {
				log.Println(err.Error())
			} else {
				items = append(items, &data)
			}
		}
	}
	return items, nil
}
