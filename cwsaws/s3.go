/*
 * File: s3.go
 * Created Date: Friday, January 26th 2024, 9:49:36 am
 *
 * Last Modified: Thu Apr 11 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsaws

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"sync"
	"time"

	"github.com/codeworks-tw/cwsutil/cwsbase"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Proxy struct {
	*s3.Client
	Context    context.Context
	BucketName string
	Versioning *bool
	Lock       sync.Mutex
}

type S3ProxyObject struct {
	Content          any
	VersionId        string
	ExpiredTimeStamp int
	NextVersionCheck int
}

var s3Objects map[string]*S3ProxyObject = map[string]*S3ProxyObject{}

func GetS3Proxy(ctx context.Context, bucketName string) S3Proxy {
	if ctx == nil {
		ctx = context.TODO()
	}

	return S3Proxy{
		Client: GetSingletonClient(ClientName_S3, ctx, func(cfg aws.Config) *s3.Client {
			return s3.NewFromConfig(cfg)
		}),
		Context:    ctx,
		BucketName: bucketName,
	}
}

func (p *S3Proxy) IsVersioning() (bool, error) {
	if p.Versioning == nil {
		o, e := p.GetBucketVersioning(p.Context, &s3.GetBucketVersioningInput{
			Bucket: aws.String(p.BucketName),
		})

		if e != nil {
			return false, e
		}
		b := o.Status == types.BucketVersioningStatusEnabled
		p.Versioning = &b
	}

	return *p.Versioning, nil
}

func (p *S3Proxy) ProxyObjectExists(subPath string) bool {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, e := p.HeadObject(p.Context, &s3.HeadObjectInput{
		Bucket: aws.String(p.BucketName),
		Key:    aws.String(subPath),
	})

	return e == nil
}

func (p *S3Proxy) ProxyPutObject(subPath string, object []byte) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, e := p.PutObject(p.Context, &s3.PutObjectInput{
		Bucket: aws.String(p.BucketName),
		Key:    aws.String(subPath),
		Body:   bytes.NewReader(object),
	})

	if e != nil {
		return e
	}

	return nil
}

func (p *S3Proxy) ProxyGetObjectVersionId(subPath string) (string, error) {
	key := p.BucketName + "/" + subPath
	if val, ok := s3Objects[key]; ok {
		return val.VersionId, nil
	}
	return "", errors.New("No avaliable s3 object: " + key)
}

func (p *S3Proxy) ProxyGetObject(subPath string, parsingFunc func(content []byte) (any, error)) (any, error) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	key := p.BucketName + "/" + subPath
	is, e := p.IsVersioning()
	if e != nil {
		return nil, e
	}
	if v, ok := s3Objects[key]; ok {
		if is {
			if int(time.Now().UTC().Unix()) < v.NextVersionCheck {
				return v.Content, nil
			}

			r, e := p.GetObjectAttributes(p.Context, &s3.GetObjectAttributesInput{
				Bucket:           aws.String(p.BucketName),
				Key:              aws.String(subPath),
				ObjectAttributes: []types.ObjectAttributes{types.ObjectAttributesChecksum},
			})

			if e != nil {
				return nil, e
			}

			if *r.VersionId == v.VersionId {
				v.updateVersionCheck()
				return v.Content, nil
			}
		} else {
			if int(time.Now().UTC().Unix()) < v.ExpiredTimeStamp {
				return v.Content, nil
			}
		}
	}

	// not exist in cache
	pb, versonId, err := p.proxyLoadObject(subPath, parsingFunc)
	if err == nil {
		s3Objects[key] = &S3ProxyObject{
			Content:          pb,
			VersionId:        versonId,
			ExpiredTimeStamp: int((time.Now().UTC().Add(time.Minute * time.Duration(cwsbase.GetEnv("S3CacheTTL", 10)))).Unix()),
		}
		s3Objects[key].updateVersionCheck()

		return pb, nil
	}

	log.Println(err)
	if v, ok := s3Objects[key]; ok {
		return v.Content, nil
	}

	return nil, err
}

func (p *S3Proxy) proxyLoadObject(subPath string, parsingFunc func(content []byte) (any, error)) (any, string, error) {
	r, e := p.GetObject(p.Context, &s3.GetObjectInput{
		Bucket: aws.String(p.BucketName),
		Key:    aws.String(subPath),
	})

	if e != nil {
		return nil, "", e
	}

	version := ""
	if p.Versioning != nil && *p.Versioning {
		version = *r.VersionId
	}

	b, e := io.ReadAll(r.Body)
	defer r.Body.Close()

	if e != nil {
		return nil, version, e
	}
	pb, e := parsingFunc(b)
	if e != nil {
		return nil, version, e
	}
	return pb, version, nil
}

func (p *S3ProxyObject) updateVersionCheck() {
	p.NextVersionCheck = int((time.Now().UTC().Add(time.Second * time.Duration(cwsbase.GetEnv("S3VersionCheck", 30)))).Unix())
}
