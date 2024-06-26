/*
 * File: sts.go
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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type STSProxy struct {
	*sts.Client
	Context context.Context
	account string
	arn     string
	userId  string
}

func GetStsClient(ctx context.Context) STSProxy {
	if ctx == nil {
		ctx = context.TODO()
	}

	return STSProxy{
		Client: GetSingletonClient[*sts.Client](ClientName_STS, ctx, func(cfg aws.Config) *sts.Client {
			return sts.NewFromConfig(cfg)
		}),
		Context: ctx,
	}
}

func (proxy *STSProxy) GetAccountId() (string, error) {
	err := proxy.initCallerIdentity()
	if err != nil {
		return "", err
	}
	return proxy.account, nil
}

func (proxy *STSProxy) initCallerIdentity() error {
	stsOut, err := proxy.GetCallerIdentity(proxy.Context, &sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}

	proxy.account = *stsOut.Account
	proxy.arn = *stsOut.Arn
	proxy.userId = *stsOut.UserId
	return nil
}
