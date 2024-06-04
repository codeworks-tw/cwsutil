/*
 * File: sqs.go
 * Created Date: Thursday, April 11th 2024, 10:31:37 am
 *
 * Last Modified: Tue Jun 04 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsaws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSProxy struct {
	*sqs.Client
	Context      context.Context
	QueueName    string
	AccountId    string
	QueueArn     string
	QueueUrl     string
	Region       string
	msgProcessFn []func(ctx context.Context, msg types.Message) error
}

func GetSqsProxy(ctx context.Context, processFn ...func(ctx context.Context, msg types.Message) error) SQSProxy {
	if ctx == nil {
		ctx = context.TODO()
	}

	var region string
	return SQSProxy{
		Client: GetSingletonClient[*sqs.Client](ClientName_SQS, ctx, func(cfg aws.Config) *sqs.Client {
			region = cfg.Region
			return sqs.NewFromConfig(cfg)
		}),
		Context:      ctx,
		Region:       region,
		msgProcessFn: processFn,
	}
}

func (proxy *SQSProxy) ProxyInitializeByName(name string) error {
	stsClient := GetStsClient(proxy.Context)
	accountId, err := stsClient.GetAccountId()
	if err != nil {
		return err
	}

	proxy.AccountId = accountId
	proxy.QueueName = name
	proxy.QueueArn = "arn:aws:sqs:" + proxy.Region + ":" + proxy.AccountId + ":" + proxy.QueueName

	out, err := proxy.GetQueueUrl(proxy.Context, &sqs.GetQueueUrlInput{
		QueueName:              &proxy.QueueName,
		QueueOwnerAWSAccountId: &proxy.AccountId,
	})
	if err != nil {
		return err
	}
	proxy.QueueUrl = *out.QueueUrl

	return nil
}

func (proxy *SQSProxy) ProxyInitializeByArn(arn string) error {
	data := strings.Split(arn, ":")
	if len(data) < 6 || data[0] != "arn" {
		return errors.New("Invalid arn: " + arn)
	}

	proxy.QueueName = data[len(data)-1]
	proxy.AccountId = data[len(data)-2]
	proxy.QueueArn = arn

	out, err := proxy.GetQueueUrl(proxy.Context, &sqs.GetQueueUrlInput{
		QueueName:              &proxy.QueueName,
		QueueOwnerAWSAccountId: &proxy.AccountId,
	})

	if err != nil {
		return err
	}

	proxy.QueueUrl = *out.QueueUrl
	return nil
}

func (proxy *SQSProxy) ProxyProcessNextMessages(MaxNumberOfMessages int32) error {
	out, err := proxy.ReceiveMessage(proxy.Context, &sqs.ReceiveMessageInput{
		QueueUrl:            &proxy.QueueUrl,
		MaxNumberOfMessages: MaxNumberOfMessages,
	})

	if err != nil {
		return err
	}

	for _, msg := range (*out).Messages {
		if len(proxy.msgProcessFn) > 0 {
			err := proxy.msgProcessFn[0](proxy.Context, msg)
			if err != nil {
				return err
			}
		}
		proxy.DeleteMessage(proxy.Context, &sqs.DeleteMessageInput{
			QueueUrl:      &proxy.QueueUrl,
			ReceiptHandle: msg.ReceiptHandle,
		})
	}
	return nil
}

func ProxyProcessLambdaEvent(ctx context.Context, event events.SQSEvent, processFn func(ctx context.Context, msg types.Message) error) error {
	for _, record := range event.Records {
		proxy := GetSqsProxy(ctx)
		err := proxy.ProxyInitializeByArn(record.EventSourceARN)
		if err != nil {
			return err
		}

		if processFn != nil {
			err := processFn(ctx, types.Message{
				Attributes:             record.Attributes,
				Body:                   &record.Body,
				MD5OfBody:              &record.Md5OfBody,
				MD5OfMessageAttributes: &record.Md5OfMessageAttributes,
				MessageAttributes:      eventMsgAttrToTypesMsgAttr(record.MessageAttributes),
				MessageId:              &record.MessageId,
				ReceiptHandle:          &record.ReceiptHandle,
			})
			if err != nil {
				return err
			}
		}
		proxy.DeleteMessage(proxy.Context, &sqs.DeleteMessageInput{
			QueueUrl:      &proxy.QueueUrl,
			ReceiptHandle: &record.ReceiptHandle,
		})

	}
	return nil
}

func eventMsgAttrToTypesMsgAttr(e map[string]events.SQSMessageAttribute) map[string]types.MessageAttributeValue {
	t := map[string]types.MessageAttributeValue{}
	for k, v := range e {
		t[k] = types.MessageAttributeValue{
			DataType:         &v.DataType,
			StringValue:      v.StringValue,
			BinaryValue:      v.BinaryValue,
			StringListValues: v.StringListValues,
			BinaryListValues: v.BinaryListValues,
		}
	}
	return t
}
