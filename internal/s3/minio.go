package s3

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	Client *minio.Client
}

func NewMinIOClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*MinIOClient, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinIOClient{Client: minioClient}, nil
}

func (m *MinIOClient) GetClient() *minio.Client {
	return m.Client
}

func (m *MinIOClient) UploadFile(bucketName, objectName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = m.Client.PutObject(context.Background(), bucketName, objectName, file, -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (m *MinIOClient) DownloadFile(bucketName, objectName, filePath string) error {
	err := m.Client.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (m *MinIOClient) UploadDirectory(bucketName, objectPrefix, dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			objectName := filepath.Join(objectPrefix, path[len(dirPath)+1:])
			_, err = m.Client.PutObject(context.Background(), bucketName, objectName, file, -1, minio.PutObjectOptions{})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (m *MinIOClient) DownloadDirectory(bucketName, objectPrefix, dirPath string) error {
	objectCh := m.Client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    objectPrefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}

		filePath := filepath.Join(dirPath, object.Key[len(objectPrefix):])
		err := m.DownloadFile(bucketName, object.Key, filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MinIOClient) ObjectExists(bucketName, objectName string) (bool, error) {
	_, err := m.Client.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if err.(minio.ErrorResponse).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (m *MinIOClient) ListFilesInDirectory(bucketName, objectPrefix string) ([]string, error) {
	var files []string

	objectCh := m.Client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    objectPrefix,
		Recursive: false,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		if !strings.HasSuffix(object.Key, "/") {
			files = append(files, object.Key)
		}
	}

	return files, nil
}
