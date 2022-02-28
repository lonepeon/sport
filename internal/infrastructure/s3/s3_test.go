package s3_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/infrastructure/s3"
	"github.com/ory/dockertest/v3"
)

const (
	bucketName     = "my.test.bucket"
	bucketUser     = "minio"
	bucketPassword = "minio123"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Run("StoreAssetFileSuccess", testStoreAssetFileSuccess)
	t.Run("StoreAssetFileInvalidCredentials", testStoreAssetFileInvalidCredentials)
	t.Run("DeleteAssetFileSuccess", testDeleteAssetFileSuccess)
	t.Run("DeleteAssetFileNotExistingFile", testDeleteAssetFileNotExistingFile)
}

func testStoreAssetFileInvalidCredentials(t *testing.T) {
	bucketEndpoint := setupS3(t)
	os.Setenv("AWS_ACCESS_KEY_ID", "NOPE")
	os.Setenv("AWS_SECRET_ACCESS_KEY", bucketPassword)
	bucket := s3.NewBucket(bucketName, "eu-west-3")
	bucket.Endpoint = bucketEndpoint

	expectedFileContent := "an important note"

	err := bucket.StoreAsset(strings.NewReader(expectedFileContent), "/a/nice/file.txt")

	testutils.AssertHasError(t, err, "shouldn't have store the file")
	testutils.AssertContainsString(t, "InvalidAccessKeyId", err.Error(), "unexpected failure")
}

func testStoreAssetFileSuccess(t *testing.T) {
	bucketEndpoint := setupS3(t)
	os.Setenv("AWS_ACCESS_KEY_ID", bucketUser)
	os.Setenv("AWS_SECRET_ACCESS_KEY", bucketPassword)
	bucket := s3.NewBucket(bucketName, "eu-west-3")
	bucket.Endpoint = bucketEndpoint

	expectedFileContent := "an important note"

	err := bucket.StoreAsset(strings.NewReader(expectedFileContent), "/a/nice/file.txt")
	testutils.AssertNoError(t, err, "can't store file")

	actualFileContent, httpStatusCode := getFile(t, bucketEndpoint, bucketName, "/a/nice/file.txt")

	testutils.AssertEqualInt(t, http.StatusOK, httpStatusCode, "unexpected response content: %v", actualFileContent)
	testutils.AssertEqualString(t, expectedFileContent, actualFileContent, "unexpected file content")
}

func testDeleteAssetFileSuccess(t *testing.T) {
	bucketEndpoint := setupS3(t)
	os.Setenv("AWS_ACCESS_KEY_ID", bucketUser)
	os.Setenv("AWS_SECRET_ACCESS_KEY", bucketPassword)
	bucket := s3.NewBucket(bucketName, "eu-west-3")
	bucket.Endpoint = bucketEndpoint

	expectedFileContent := "an important note"

	err := bucket.StoreAsset(strings.NewReader(expectedFileContent), "/a/nice/file.txt")
	testutils.AssertNoError(t, err, "should have store the file")

	actualFileContent, httpStatusCode := getFile(t, bucketEndpoint, bucketName, "/a/nice/file.txt")
	testutils.AssertEqualInt(t, http.StatusOK, httpStatusCode, "unexpected response content: %v", actualFileContent)

	err = bucket.DeleteAsset("/a/nice/file.txt")
	testutils.AssertNoError(t, err, "should have deleted file")

	actualFileContent, httpStatusCode = getFile(t, bucketEndpoint, bucketName, "/a/nice/file.txt")
	testutils.AssertEqualInt(t, http.StatusNotFound, httpStatusCode, "unexpected response content: %v", actualFileContent)
}

func testDeleteAssetFileNotExistingFile(t *testing.T) {
	bucketEndpoint := setupS3(t)
	os.Setenv("AWS_ACCESS_KEY_ID", bucketUser)
	os.Setenv("AWS_SECRET_ACCESS_KEY", bucketPassword)
	bucket := s3.NewBucket(bucketName, "eu-west-3")
	bucket.Endpoint = bucketEndpoint

	err := bucket.DeleteAsset("/a/non-existing/file.txt")
	testutils.AssertNoError(t, err, "should have not failed when deleting a non-existing the file")
}

func getFile(t *testing.T, bucketEndpoint, bucketName, filePath string) (string, int) {
	resp, err := http.Get(bucketEndpoint + "/" + bucketName + filePath)
	testutils.AssertNoError(t, err, "can't get back stored file")

	body, err := io.ReadAll(resp.Body)
	testutils.AssertNoError(t, err, "can't read received file")

	content := string(body)

	return content, resp.StatusCode
}

func setupS3(t *testing.T) string {
	pool, err := dockertest.NewPool("")
	testutils.AssertNoError(t, err, "can't connnect to docker")

	minio, minioEndpoint := startMinioContainer(t, pool, bucketUser, bucketPassword)

	// this is just in case the purge doesn't work 😅
	testutils.AssertNoError(t, minio.Expire(60), "can't set automatic purge")

	mc, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "minio/mc",
		Links:      []string{minio.Container.Name},
		Entrypoint: []string{"/bin/sh", "-c", fmt.Sprintf(`
			/usr/bin/mc alias set minio http:/%[1]s:9000 %[2]s %[3]s;
			/usr/bin/mc mb minio/%[4]s;
			/usr/bin/mc policy set public minio/%[4]s;
			exit 0;`,
			minio.Container.Name, bucketUser, bucketPassword, bucketName)},
	})
	testutils.AssertNoError(t, err, "can't start minio configuration container")

	// this is just in case the purge doesn't work 😅
	testutils.AssertNoError(t, minio.Expire(30), "can't set automatic purge")

	var retries = 10
	for retries > 0 && mc.Container.State.FinishedAt.IsZero() {
		time.Sleep(1 * time.Second)
		retries--
	}

	t.Cleanup(func() {
		testutils.AssertNoError(t, pool.Purge(minio), "can't remove minio container")
		testutils.AssertNoError(t, pool.Purge(mc), "can't remove minio configuration container")
	})

	return minioEndpoint
}

func startMinioContainer(t *testing.T, pool *dockertest.Pool, bucketUser, bucketPassword string) (*dockertest.Resource, string) {
	minio, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "minio/minio",
		Tag:        "RELEASE.2022-01-08T03-11-54Z",
		Env: []string{
			fmt.Sprintf("MINIO_ROOT_USER=%s", bucketUser),
			fmt.Sprintf("MINIO_ROOT_PASSWORD=%s", bucketPassword),
		},
		ExposedPorts: []string{"9000/tcp"},
		Cmd:          []string{"server", "/data"},
	})
	testutils.AssertNoError(t, err, "can't start minio container")

	minioEndpoint := fmt.Sprintf("http://localhost:%s", minio.GetPort("9000/tcp"))

	err = pool.Retry(func() error {
		url := fmt.Sprintf("%s/minio/health/live", minioEndpoint)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code not OK")
		}
		testutils.AssertNoError(t, err, "Could not connect to minio container")

		return nil
	})

	return minio, minioEndpoint
}
