package s3

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	Client *s3.S3
}

func NewS3Client() (*S3Client, error) {
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("BUCKET_ENDPOINT")
	useSSL := os.Getenv("USE_SSL") == "true"

	sessOpts := session.Options{
		Config: aws.Config{
			Region:           aws.String(region),
			S3ForcePathStyle: aws.Bool(true),
		},
	}

	if endpoint != "" {
		fmt.Println("Configuring S3 client for local development")
		sessOpts.Config.Endpoint = aws.String(endpoint)
		sessOpts.Config.DisableSSL = aws.Bool(!useSSL)
		sessOpts.Config.Credentials = credentials.NewStaticCredentials(
			os.Getenv("BUCKET_ACCESS_KEY_ID"),
			os.Getenv("BUCKET_SECRET_ACCESS_KEY"),
			"",
		)
	} else {
		fmt.Println("Configuring S3 client for AWS deployment")
		sessOpts.Config.Region = aws.String(region)
		if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
			sessOpts.Config.Credentials = credentials.NewStaticCredentials(
				os.Getenv("BUCKET_ACCESS_KEY_ID"),
				os.Getenv("BUCKET_SECRET_ACCESS_KEY"),
				"",
			)
		}
	}

	sess, err := session.NewSessionWithOptions(sessOpts)
	if err != nil {
		return nil, err
	}

	return &S3Client{Client: s3.New(sess)}, nil
}
func (s *S3Client) GetClient() *s3.S3 {
	return s.Client
}

func (s *S3Client) CreateBucket(bucketName string) error {
	_, err := s.Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeBucketAlreadyOwnedByYou {
			return nil
		}
		return err
	}
	return nil
}

func (s *S3Client) BucketExists(bucketName string) (bool, error) {
	_, err := s.Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket, "NotFound", "Forbidden":
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func (s *S3Client) UploadFile(bucketName, objectName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
		Body:   file,
	})
	return err
}

func (s *S3Client) DownloadFile(bucketName, objectName, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Client) UploadDirectory(bucketName, objectPrefix, dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return err
			}
			objectKey := filepath.Join(objectPrefix, relativePath)
			return s.UploadFile(bucketName, objectKey, path)
		}
		return nil
	})
}

func (s *S3Client) DownloadDirectory(bucketName, objectPrefix, dirPath string) error {
	resp, err := s.Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(objectPrefix),
	})
	if err != nil {
		return err
	}

	for _, item := range resp.Contents {
		filePath := filepath.Join(dirPath, strings.TrimPrefix(*item.Key, objectPrefix))
		if err := s.DownloadFile(bucketName, *item.Key, filePath); err != nil {
			return err
		}
	}

	return nil
}

func (s *S3Client) ObjectExists(bucketName, objectName string) (bool, error) {
	_, err := s.Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *S3Client) ListFilesInDirectory(bucketName, objectPrefix string) ([]string, error) {
	var files []string

	resp, err := s.Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(objectPrefix),
	})
	if err != nil {
		return nil, err
	}

	for _, item := range resp.Contents {
		if !strings.HasSuffix(*item.Key, "/") {
			files = append(files, *item.Key)
		}
	}

	return files, nil
}
