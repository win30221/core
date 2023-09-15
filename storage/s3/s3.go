package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Config struct {
	// Bucket 可以在 `https://s3.console.aws.amazon.com/s3/home` 登入後找到
	Bucket string
	// Region 可用的地區請參考 https://docs.aws.amazon.com/directoryservice/latest/admin-guide/regions.html
	Region          string
	AccessKeyID     string
	SecretAccessKey string

	// RootSubset 將 bucket 底下資源區隔為 dev, prd。(dev, qa, sit, pt 等站點 root_folder 都用 dev)
	RootSubset string
	CdnURL     string
}

type Storage struct {
	client   *s3.S3
	uploader *s3manager.Uploader
	conf     Config
}

func New(conf Config) (s *Storage, err error) {
	awsConfig := &aws.Config{
		Region:      aws.String(conf.Region),
		Credentials: credentials.NewStaticCredentials(conf.AccessKeyID, conf.SecretAccessKey, ""),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return
	}

	s = &Storage{
		client:   s3.New(sess),
		uploader: s3manager.NewUploader(sess),
		conf:     conf,
	}

	return
}

// UploadImage
func (s *Storage) UploadImage(ctx context.Context, dirPath, filename, contentType string, file io.ReadSeeker) (imageURL string, err error) {
	// Key 開頭不能有 `/`
	// - OK: 	"dev/image/test.png"
	// - Wrong: "/dev/image/test.png"
	filePath := fmt.Sprintf("%s/%s/%s", s.conf.RootSubset, dirPath, filename)

	_, err = s.uploader.UploadWithContext(
		ctx,
		&s3manager.UploadInput{
			Bucket:      aws.String(s.conf.Bucket),
			Key:         aws.String(filePath),
			ContentType: aws.String(contentType),
			Body:        file,
		},
	)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			err = fmt.Errorf("s3storage.PutObjectWithContext failed, code: %s, message: %s, originErr: %s", awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
		}

		return
	}

	imageURL = s.conf.CdnURL + "/" + filePath

	return
}

// DeleteImage
func (s *Storage) DeleteImage(ctx context.Context, dirPath, filename string) (err error) {
	// Key 開頭不能有 `/`
	// - OK: 	"dev/image/test.png"
	// - Wrong: "/dev/image/test.png"
	filePath := fmt.Sprintf("%s/%s/%s", s.conf.RootSubset, dirPath, filename)

	_, err = s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.conf.Bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			err = fmt.Errorf("s3storage.Delete failed, code: %s, message: %s, originErr: %s", awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
		}

		return
	}

	return
}
