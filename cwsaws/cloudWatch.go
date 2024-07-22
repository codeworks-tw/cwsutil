/*
 * File: cloudWatch.go
 * Created Date: Monday, July 22nd 2024, 5:28:16 pm
 *
 * Last Modified: Mon Jul 22 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsaws

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type CloudWatchLogsProxy struct {
	*cloudwatchlogs.Client
	Context  *context.Context
	LogGroup string
}

func GetCloudWatchLogProxy(logGroup string, ctx context.Context, optFns ...func(*config.LoadOptions) error) CloudWatchLogsProxy {
	if ctx == nil {
		ctx = context.TODO()
	}

	return CloudWatchLogsProxy{
		LogGroup: logGroup,
		Context:  &ctx,
		Client: GetSingletonClient(ClientName_CloudWatch, ctx, func(cfg aws.Config) *cloudwatchlogs.Client {
			return cloudwatchlogs.NewFromConfig(cfg)
		}, optFns...),
	}
}

func (p *CloudWatchLogsProxy) GetTimedLogStreamName(interval int) *string {
	current := time.Now().UTC()
	name := current.Format("2006-01-02") + "[" + strconv.Itoa(current.Hour()) + "-" + strconv.Itoa(current.Minute()/interval*interval) + "]"
	return &name
}

func (p *CloudWatchLogsProxy) CreateLogGroup() error {
	dr, err := p.Client.DescribeLogGroups(*p.Context, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &p.LogGroup,
		Limit:              aws.Int32(1),
	})
	if err != nil {
		return err
	}

	if len(dr.LogGroups) == 0 {
		_, err := p.Client.CreateLogGroup(*p.Context, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: &p.LogGroup,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *CloudWatchLogsProxy) CreateLogStream() error {
	err := p.CreateLogGroup()
	if err != nil {
		return err
	}

	dr, err := p.Client.DescribeLogStreams(*p.Context, &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        &p.LogGroup,
		LogStreamNamePrefix: p.GetTimedLogStreamName(15),
		Limit:               aws.Int32(1),
	})
	if err != nil {
		return err
	}

	if len(dr.LogStreams) == 0 {
		_, err := p.Client.CreateLogStream(*p.Context, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  &p.LogGroup,
			LogStreamName: p.GetTimedLogStreamName(15),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *CloudWatchLogsProxy) SendMessage(message string) error {
	p.CreateLogStream()
	_, err := p.Client.PutLogEvents(*p.Context, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  &p.LogGroup,
		LogStreamName: p.GetTimedLogStreamName(15),
		LogEvents: []types.InputLogEvent{
			{
				Message:   &message,
				Timestamp: aws.Int64(int64(time.Now().UnixMilli())),
			},
		},
	})
	return err
}
