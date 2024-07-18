package s3

import (
	"fmt"
	"io"
	"net/http"
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

func NewS3Client(checkpoint ...bool) (*S3Client, error) {
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("BUCKET_ENDPOINT")
	useSSL := os.Getenv("USE_SSL") == "true"

	var sess *session.Session
	var err error

	var presignedURLEndpoint string

	forCheckpoint := false
	if len(checkpoint) > 0 {
		forCheckpoint = checkpoint[0]
	}
	//below change only for checkpoints
	if forCheckpoint && endpoint == "http://object-store:9000" {
		presignedURLEndpoint = "http://localhost:9000"
	} else {
		presignedURLEndpoint = endpoint
	}

	if endpoint != "" {
		fmt.Println("Configuring S3 client for local development")
		sessOpts := session.Options{
			Config: aws.Config{
				Region:           aws.String(region),
				Endpoint:         aws.String(presignedURLEndpoint),
				S3ForcePathStyle: aws.Bool(true),
				DisableSSL:       aws.Bool(!useSSL),
				Credentials: credentials.NewStaticCredentials(
					os.Getenv("BUCKET_ACCESS_KEY_ID"),
					os.Getenv("BUCKET_SECRET_ACCESS_KEY"),
					"",
				),
			},
		}
		sess, err = session.NewSessionWithOptions(sessOpts)
	} else {
		fmt.Println("Configuring S3 client for AWS deployment")
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
		if err != nil {
			fmt.Println("Error creating session for S3 client:")
			fmt.Println(err)
			return nil, err
		}
	}

	if err != nil {
		fmt.Println("Error creating session for S3 client:")
		fmt.Println(err)
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

func (s *S3Client) DownloadFile(bucketName, objectName, fileName string) error {
	// Create a new file in the provided path.
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the object from S3 and write its content to the file.
	output, err := s.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return err
	}
	defer output.Body.Close()

	// Copy data from S3 object to the file
	_, err = io.Copy(file, output.Body)
	return err
}

func (s *S3Client) StreamFileToResponse(s3URI string, w http.ResponseWriter, filename string) error {
	bucketName, objectName, err := s.GetBucketAndKeyFromURI(s3URI)
	if err != nil {
		return err
	}
	output, err := s.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return err
	}
	defer output.Body.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err = io.Copy(w, output.Body); err != nil {
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

func (s *S3Client) GetBucketAndKeyFromURI(uri string) (string, string, error) {
	uriParts := strings.Split(uri, "://")
	if len(uriParts) != 2 {
		return "", "", fmt.Errorf("invalid URI: %s", uri)
	}
	uriParts = strings.Split(uriParts[1], "/")
	bucket := uriParts[0]
	path := strings.Join(uriParts[1:], "/")
	return bucket, path, nil
}
