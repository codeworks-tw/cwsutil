/*
 * File: ses.go
 * Created Date: Tuesday, September 26th 2023, 2:05:51 pm
 *
 * Last Modified: Wed Oct 18 2023
 * Modified By: Howard Ling-Hao Kung
 */

package awsutil

import (
	"context"
	"net/mail"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SESProxy struct {
	*ses.Client
	Context context.Context
}

type SESProxyEmailInput struct {
	From           mail.Address
	To             []mail.Address
	CC             []mail.Address
	BCC            []mail.Address
	Subject        string
	SubjectCharset string // default: UTF-8
	Body           string
	BodyCharset    string // default: UTF-8
	IsText         bool   // default: html
}

type SESProxyTemplateInput struct {
	TemplateName  string
	SubjectParams map[string]any
	MessageParams map[string]any
	IsText        bool // default: html
}

type SESProxyTemplateEmailInput struct {
	TemplateInput SESProxyTemplateInput
	From          mail.Address
	To            []mail.Address
	CC            []mail.Address
	BCC           []mail.Address
}

func GetSesProxy(ctx context.Context) SESProxy {
	if ctx == nil {
		ctx = context.TODO()
	}

	return SESProxy{
		Client: GetSingletonClient(ClientName_SES, ctx, func(cfg aws.Config) *ses.Client {
			return ses.NewFromConfig(cfg)
		}),
		Context: ctx,
	}
}

func (p *SESProxy) ProxySendTemplateEmail(input *SESProxyTemplateEmailInput) (*ses.SendEmailOutput, error) {
	subject, message, e := p.ProxyGetTemplateMessage(input.TemplateInput)
	if e != nil {
		return nil, e
	}

	emailInput := &SESProxyEmailInput{
		From: input.From,
		To:   input.To,
		CC:   input.CC,
		BCC:  input.BCC,
	}

	emailInput.Subject = subject
	emailInput.Body = message
	emailInput.IsText = input.TemplateInput.IsText
	return p.ProxySendEmail(emailInput)
}

func (p *SESProxy) ProxySendEmail(input *SESProxyEmailInput) (*ses.SendEmailOutput, error) {
	data, e := processInput(input)
	if e != nil {
		return nil, e
	}
	return p.SendEmail(p.Context, data)
}

func (p *SESProxy) ProxyGetTemplateMessage(input SESProxyTemplateInput) (string, string, error) {
	o, e := p.GetTemplate(p.Context, &ses.GetTemplateInput{
		TemplateName: &input.TemplateName,
	})
	if e != nil {
		return "", "", e
	}

	subject := *o.Template.SubjectPart
	for k, v := range input.SubjectParams {
		subject = strings.ReplaceAll(subject, k, v.(string))
	}

	if input.IsText {
		message := *o.Template.TextPart
		for k, v := range input.MessageParams {
			message = strings.ReplaceAll(message, k, v.(string))
		}
		return subject, message, nil
	}

	message := *o.Template.HtmlPart
	for k, v := range input.MessageParams {
		message = strings.ReplaceAll(message, k, v.(string))
	}
	return subject, message, nil
}

func verifyEmail(input *SESProxyEmailInput) error {
	if _, e := mail.ParseAddress((input.From.Address)); e != nil {
		return e
	}
	for _, m := range input.To {
		if _, e := mail.ParseAddress(m.Address); e != nil {
			return e
		}
	}
	for _, m := range input.CC {
		if _, e := mail.ParseAddress(m.Address); e != nil {
			return e
		}
	}
	for _, m := range input.BCC {
		if _, e := mail.ParseAddress(m.Address); e != nil {
			return e
		}
	}
	return nil
}

func processInput(input *SESProxyEmailInput) (*ses.SendEmailInput, error) {
	dest := &types.Destination{}
	for _, m := range input.To {
		dest.ToAddresses = append(dest.ToAddresses, m.Address)
	}
	for _, m := range input.CC {
		dest.ToAddresses = append(dest.ToAddresses, m.Address)
	}
	for _, m := range input.BCC {
		dest.ToAddresses = append(dest.ToAddresses, m.Address)
	}

	subject := types.Content{
		Data: &input.Subject,
	}

	if input.SubjectCharset != "" {
		subject.Charset = &input.SubjectCharset
	} else {
		subject.Charset = aws.String("UTF-8")
	}

	body := types.Content{
		Data: &input.Body,
	}

	if input.BodyCharset != "" {
		body.Charset = &input.BodyCharset
	} else {
		body.Charset = aws.String("UTF-8")
	}

	e := verifyEmail(input)
	if e != nil {
		return nil, e
	}

	if input.IsText {
		return &ses.SendEmailInput{
			Source:      &input.From.Address,
			Destination: dest,
			Message: &types.Message{
				Subject: &subject,
				Body: &types.Body{
					Text: &body,
				},
			},
		}, nil
	}
	return &ses.SendEmailInput{
		Source:      &input.From.Address,
		Destination: dest,
		Message: &types.Message{
			Subject: &subject,
			Body: &types.Body{
				Html: &body,
			},
		},
	}, nil
}
