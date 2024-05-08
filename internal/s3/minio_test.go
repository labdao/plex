package s3

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	minioClient *MinIOClient
	bucketName  string
	tempDir     string
)

func setup(t *testing.T) {
	endpoint := ""
	accessKeyID := ""
	secretAccessKey := ""
	useSSL := false

	var err error
	minioClient, err = NewMinIOClient(endpoint, accessKeyID, secretAccessKey, useSSL)
	require.NoError(t, err)

	bucketName = "test-bucket"
	err = minioClient.Client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	require.NoError(t, err)

	tempDir, err = os.MkdirTemp("", "minio-test")
	require.NoError(t, err)
}

func teardown(t *testing.T) {
	err := minioClient.Client.RemoveBucket(context.Background(), bucketName)
	require.NoError(t, err)

	err = os.RemoveAll(tempDir)
	assert.NoError(t, err)
}

func TestUploadAndDownloadFile(t *testing.T) {
	setup(t)
	defer teardown(t)

	testFilePath := filepath.Join(tempDir, "test-file.txt")
	err := os.WriteFile(testFilePath, []byte("Hello, world!"), 0644)
	require.NoError(t, err)

	objectName := "test-object.txt"
	err = minioClient.UploadFile(bucketName, objectName, filepath.Join(tempDir, objectName))
	assert.NoError(t, err)

	exists, err := minioClient.ObjectExists(bucketName, "downloaded-file.txt")
	assert.NoError(t, err)
	assert.True(t, exists)

	downloadFilePath := filepath.Join(tempDir, "downloaded-file.txt")
	err = minioClient.DownloadFile(bucketName, objectName, downloadFilePath)
	assert.NoError(t, err)

	content, err := os.ReadFile(downloadFilePath)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", string(content))
}

func TestUploadAndDownloadDirectory(t *testing.T) {
	setup(t)
	defer teardown(t)

	testDirPath := filepath.Join(tempDir, "test-dir")
	err := os.Mkdir(testDirPath, 0755)
	require.NoError(t, err)

	testTextFilePath := filepath.Join(testDirPath, "test-file.txt")
	err = os.WriteFile(testTextFilePath, []byte("Hello, world from .txt!"), 0644)
	require.NoError(t, err)

	testJSONFilePath := filepath.Join(testDirPath, "test-file.json")
	err = os.WriteFile(testJSONFilePath, []byte(`{"message": "Hello, world from .json!"}`), 0644)
	require.NoError(t, err)

	objectPrefix := "test-dir/"
	err = minioClient.UploadDirectory(bucketName, objectPrefix, testDirPath)
	assert.NoError(t, err)

	files, err := minioClient.ListFilesInDirectory(bucketName, objectPrefix)
	assert.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Contains(t, files, "test-file.txt")
	assert.Contains(t, files, "test-file.json")

	downloadDirPath := filepath.Join(tempDir, "downloaded-dir")
	err = minioClient.DownloadDirectory(bucketName, objectPrefix, downloadDirPath)
	assert.NoError(t, err)

	downloadedTextFilePath := filepath.Join(downloadDirPath, "test-file.txt")
	assert.FileExists(t, downloadedTextFilePath)

	downloadedJSONFilePath := filepath.Join(downloadDirPath, "test-file.json")
	assert.FileExists(t, downloadedJSONFilePath)

	textContent, err := os.ReadFile(downloadedTextFilePath)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world from .txt!", string(textContent))

	jsonContent, err := os.ReadFile(downloadedJSONFilePath)
	assert.NoError(t, err)
	assert.Equal(t, `{"message": "Hello, world from .json!"}`, string(jsonContent))
}
