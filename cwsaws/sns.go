/*
 * File: sns.go
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
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSProxy struct {
	*sns.Client
	Context  context.Context
	Topic    string
	TopicArn string
}

type SNSProxyTemplateNotificationInput struct {
	TemplateName string
	Phone        string
	Params       []any
}

func GetSnsProxy(ctx context.Context) SNSProxy {
	return SNSProxy{
		Client: GetSingletonClient(ClientName_SNS, ctx, func(cfg aws.Config) *sns.Client {
			return sns.NewFromConfig(cfg)
		}),
		Context: ctx,
	}
}

func GetSnsProxyWithTopic(ctx context.Context, Topic string) SNSProxy {
	return SNSProxy{
		Client: GetSingletonClient(ClientName_SNS, ctx, func(cfg aws.Config) *sns.Client {
			return sns.NewFromConfig(cfg)
		}),
		Context: ctx,
		Topic:   Topic,
	}
}

func (s *SNSProxy) ProxySendPhoneNotification(phone string, content string, phases ...any) (*sns.PublishOutput, error) {
	o, e := s.Publish(s.Context, &sns.PublishInput{
		Message:     aws.String(fmt.Sprintf(content, phases...)),
		PhoneNumber: aws.String(phone),
	})
	if e != nil {
		return nil, e
	}
	return o, nil
}

func (s *SNSProxy) ProxySendTemplateNotification(input *SNSProxyTemplateNotificationInput) (*sns.PublishOutput, error) {
	sesClient := GetSesProxy(s.Context)

	template, err := sesClient.GetTemplate(s.Context, &ses.GetTemplateInput{
		TemplateName: aws.String(input.TemplateName),
	})
	if err != nil {
		return nil, err
	}
	return s.ProxySendPhoneNotification(input.Phone, *template.Template.TextPart, input.Params...)
}
