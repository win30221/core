package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/win30221/core/http/catch"
	"github.com/win30221/core/syserrno"
)

type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	// Bucket 可以在 `https://s3.console.aws.amazon.com/s3/home` 登入後找到
	Bucket string
	// Region 可用的地區請參考 https://docs.aws.amazon.com/directoryservice/latest/admin-guide/regions.html
	Region          string
	EndPointVersion int
	Endpoint        string
	PublicUrl       string
	ACL             string

	// RootSubset 將 bucket 底下資源區隔為 dev, prd。(dev, qa, sit, pt 等站點 root_folder 都用 dev)
	RootSubset string
}

type Storage struct {
	client *s3.Client
	conf   Config
}

type resolverV2 struct {
	endPoint string
}

func (r *resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (res smithyendpoints.Endpoint, err error) {
	u, err := url.Parse(r.endPoint)
	if err != nil {
		return
	}

	res = smithyendpoints.Endpoint{
		URI: *u,
	}
	return
}

func New(conf Config) (s *Storage, err error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(conf.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.AccessKeyID, conf.SecretAccessKey, "")),
	)
	if err != nil {
		return
	}

	var client *s3.Client
	switch conf.EndPointVersion {
	case 1: // cloudflare
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(conf.Endpoint)
		})
	case 2: // linode
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.EndpointResolverV2 = &resolverV2{conf.Endpoint}
		})
	}

	s = &Storage{
		client: client,
		conf:   conf,
	}

	return
}

// UploadImage
func (s *Storage) UploadImage(ctx context.Context, dirPath, filename, contentType string, file multipart.File) (imageURL string, err error) {
	// Key 開頭不能有 `/`
	// - OK: 	"dev/image/test.png"
	// - Wrong: "/dev/image/test.png"
	filePath := fmt.Sprintf("%s/%s/%s", s.conf.RootSubset, dirPath, filename)

	_, err = s.client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(s.conf.Bucket),
			Key:         aws.String(filePath),
			ContentType: aws.String(contentType),
			Body:        file,
			ACL:         types.ObjectCannedACL(s.conf.ACL),
		},
	)
	if err != nil {
		err = catch.NewWitStack(syserrno.AWSS3, "s3.Upload failed", fmt.Sprintf("s3.Upload failed. err: %s", err.Error()), 3)
		return
	}

	imageURL = s.conf.PublicUrl + "/" + filePath

	return
}

// DeleteImage
func (s *Storage) DeleteImage(ctx context.Context, dirPath, filename string) (err error) {
	// Key 開頭不能有 `/`
	// - OK: 	"dev/image/test.png"
	// - Wrong: "/dev/image/test.png"
	filePath := fmt.Sprintf("%s/%s/%s", s.conf.RootSubset, dirPath, filename)

	_, err = s.client.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.conf.Bucket),
			Key:    aws.String(filePath),
		},
	)
	if err != nil {
		err = catch.NewWitStack(syserrno.AWSS3, "s3.Upload failed", fmt.Sprintf("s3.Upload failed. err: %s", err.Error()), 3)
		return
	}

	return
}
