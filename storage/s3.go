package storage

import (
	"log"

	"github.com/win30221/core/config"
	"github.com/win30221/core/storage/s3"
)

func GetS3(path string) (s *s3.Storage) {
	var err error

	defer func() {
		if err != nil {
			log.Fatalf("get aws s3 error: %s \n - path %s \n", err, path)
		}
	}()

	accessKeyID, _ := config.GetString(path+"/access_key_id", true)
	secretAccessKey, _ := config.GetString(path+"/secret_access_key", true)
	bucket, _ := config.GetString(path+"/bucket", true)
	region, _ := config.GetString(path+"/region", true)
	endpoint, _ := config.GetString(path+"/endpoint", true)
	acl, _ := config.GetString(path+"/acl", true)
	rootSubset, _ := config.GetString(path+"/root_subset", true)
	cdnURL, _ := config.GetString(path+"/cdn_url", true)

	s, err = s3.New(s3.Config{
		Bucket:          bucket,
		Region:          region,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Endpoint:        endpoint,
		ACL:             acl,

		RootSubset: rootSubset,
		CdnURL:     cdnURL,
	})

	return
}
