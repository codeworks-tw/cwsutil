/*
 * File: aws_clients.go
 * Created Date: Sunday, September 17th 2023, 6:49:44 pm
 *
 * Last Modified: Thu Oct 12 2023
 * Modified By: Howard Ling-Hao Kung
 */

package awsutil

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type ClientName string

var lock sync.Mutex
var clients map[ClientName]any = map[ClientName]any{}

const (
	ClientName_STS      ClientName = "STS"
	ClientName_DynamoDB ClientName = "DynamoDB"
	ClientName_SQS      ClientName = "SQS"
	ClientName_SNS      ClientName = "SNS"
	ClientName_S3       ClientName = "S3"
	ClientName_SES      ClientName = "SES"
)

func GetSingletonClient[T any](name ClientName, ctx context.Context, clientGenFn func(cfg aws.Config) T, optFns ...func(*config.LoadOptions) error) T {
	lock.Lock()
	defer lock.Unlock()

	client, ok := clients[name]

	if !ok {
		cfg, err := config.LoadDefaultConfig(ctx, optFns...)
		if err != nil {
			log.Fatalln(err.Error())
		}

		clients[name] = clientGenFn(cfg)
		client = clients[name]
	}

	return client.(T)
}
