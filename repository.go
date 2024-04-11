/*
 * File: repository.go
 * Created Date: Saturday, January 27th 2024, 9:46:26 am
 *
 * Last Modified: Thu Apr 11 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsutil

import (
	"context"
	"reflect"

	"github.com/codeworks-tw/cwsutil/cwsaws"
	"github.com/codeworks-tw/cwsutil/cwsbase"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type IRepository[PKey any] interface {
	Get(ctx context.Context, pKey PKey, columns ...string) (*map[string]any, error)
	Query(ctx context.Context, indexName string, expr expression.Expression) ([]*map[string]any, error)
	Merge(ctx context.Context, pKey PKey, expr expression.Expression) (*map[string]any, error)
	Delete(ctx context.Context, pKey PKey) (*map[string]any, error)
	GetDynamoDBTableProxy(ctx context.Context) *cwsaws.DynamoDBTableProxy[map[string]any]
	GetPKeyKeys() []string
}

type Repository[PKey any] struct {
	IRepository[PKey]
	TableName string
}

func (r *Repository[PKey]) GetPKeyKeys() []string {
	var t PKey
	val := reflect.Indirect(reflect.ValueOf(t))
	keys := []string{}
	for i := 0; i < val.Type().NumField(); i++ {
		keys = append(keys, val.Type().Field(i).Name)

	}
	return keys
}

func (r *Repository[PKey]) GetDynamoDBTableProxy(ctx context.Context) cwsaws.DynamoDBTableProxy[map[string]any] {
	if cwsbase.GetEnvironmentInfo().IsLocal {
		credential := config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cwsbase.GetEnv[string]("Local_DynamoDB_AWS_ID"),
			cwsbase.GetEnv[string]("Local_DynamoDB_AWS_Secret"), ""))
		endPoint := config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           cwsbase.GetEnv[string]("Local_DynamoDB_URL"),
					SigningRegion: cwsbase.GetEnv[string]("Local_DynamoDB_REGION"),
				}, nil
			}))

		return cwsaws.GetDynamoDBTableProxy[map[string]any](r.TableName, ctx, credential, endPoint)
	}
	return cwsaws.GetDynamoDBTableProxy[map[string]any](r.TableName, ctx)
}

func (r *Repository[PKey]) Get(ctx context.Context, pKey PKey, columns ...string) (*map[string]any, error) {
	tableProxy := r.GetDynamoDBTableProxy(ctx)

	keys, err := attributevalue.MarshalMap(pKey)
	if err != nil {
		return nil, err
	}

	// var e chan error
	if len(columns) > 0 {
		proj := expression.ProjectionBuilder{}
		for _, c := range columns {
			proj = expression.AddNames(proj, expression.Name(c))
		}
		expr, err := expression.NewBuilder().WithProjection(proj).Build()
		if err != nil {
			return nil, err
		}

		return tableProxy.ProxyGetItem(&dynamodb.GetItemInput{
			Key:                      keys,
			ProjectionExpression:     expr.Projection(),
			ExpressionAttributeNames: expr.Names(),
		})
	} else {
		return tableProxy.ProxyGetItem(&dynamodb.GetItemInput{
			Key: keys,
		})
	}
}

func (r *Repository[PKey]) Merge(ctx context.Context, pKey PKey, expr expression.Expression) (*map[string]any, error) {
	tableProxy := r.GetDynamoDBTableProxy(ctx)

	keys, err := attributevalue.MarshalMap(pKey)
	if err != nil {
		return nil, err
	}

	return tableProxy.ProxyUpdateItem(&dynamodb.UpdateItemInput{
		Key:                       keys,
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		ReturnValues:              types.ReturnValueAllNew, // always return in this case
	})
}

func (r *Repository[PKey]) Delete(ctx context.Context, pKey PKey) (*map[string]any, error) {
	tableProxy := r.GetDynamoDBTableProxy(ctx)

	keys, err := attributevalue.MarshalMap(pKey)
	if err != nil {
		return nil, err
	}

	return tableProxy.ProxyDeleteItem(&dynamodb.DeleteItemInput{
		Key:          keys,
		ReturnValues: types.ReturnValueAllOld, // always return in this case
	})
}

func (r *Repository[PKey]) Query(ctx context.Context, indexName string, expr expression.Expression) ([]*map[string]any, error) {
	tableProxy := r.GetDynamoDBTableProxy(ctx)
	return tableProxy.ProxyQuery(&dynamodb.QueryInput{
		IndexName:                 aws.String(indexName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})
}

func (r *Repository[PKey]) QueryBatchUpdate(ctx context.Context, indexName string, expr expression.Expression, callback func(item *map[string]any, batchRequests *[]types.WriteRequest) *[]types.WriteRequest) {
	tableProxy := r.GetDynamoDBTableProxy(ctx)

	w := tableProxy.ProxyQueryBatchUpdate(&dynamodb.QueryInput{
		IndexName:                 aws.String(indexName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}, callback)
	w.Wait()
}
