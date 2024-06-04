package s3

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
	// err = m.Client.FGetObject(context.Background(), "convexity", "uploaded/96942c2c5f73f89849f3ff183dafd864e350dbeaa899a7ba4fcce3c7fcaaf50d/a1bf7360-5f8f-46e2-8461-850d68e15d00_unrelaxed_rank_001_alphafold2_multimer_v3_model_3_seed_000 (1).pdb", "uploaded/96942c2c5f73f89849f3ff183dafd864e350dbeaa899a7ba4fcce3c7fcaaf50d/a1bf7360-5f8f-46e2-8461-850d68e15d00_unrelaxed_rank_001_alphafold2_multimer_v3_model_3_seed_000 (1).pdb", minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (m *MinIOClient) StreamFileToResponse(bucketName, objectName string, w http.ResponseWriter, filename string) error {
	object, err := m.Client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err = io.Copy(w, object); err != nil {
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
