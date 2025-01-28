package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	S3Client *s3.Client
}

func (bucket Bucket) uploadFile(
	ctx context.Context,
	fPath string,
	fileName string,
) error {
	pId := os.Getenv("PROJECT_ID")
	bName := os.Getenv("BUCKET_NAME")

	file, openFileErr := os.Open(fPath)
	if openFileErr != nil {
		return nil
	}

	defer file.Close()

	data, readFileErr := io.ReadAll(file)
	if readFileErr != nil {
		return nil
	}

	mimeType := http.DetectContentType(data)
	bucketName := bName
	objectKey := "_output/" + pId + "/" + fileName

	largeBuffer := bytes.NewReader(data)
	var partMiBs int64 = 10

	uploader := manager.NewUploader(bucket.S3Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        largeBuffer,
		ContentType: &mimeType,
	})
	if err != nil {
		log.Printf("Error while uploading object to %s. %s\n", bucketName, err)
	}

	fmt.Printf(
		"Successfully uploaded file %s to bucket %s with key %s\n",
		fPath,
		bucketName,
		objectKey,
	)

	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Init() {
	rootDir, getRootErr := os.Getwd()
	checkErr(getRootErr)

	buildProject(rootDir)
	uploadProject(rootDir)
}

func uploadProject(rootDir string) {
	filePaths := []string{}

	readDistErr := filepath.Walk(rootDir+"/output"+"/dist",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories, only add files
			if !info.IsDir() {
				filePaths = append(filePaths, path)
			}

			return nil
		})

	checkErr(readDistErr)

	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	checkErr(err)

	// Create an S3 service client
	s3Svc := s3.NewFromConfig(cfg)

	bucket := Bucket{s3Svc}

	// uploading files one by one
	for _, fp := range filePaths {
		fArr := strings.Split(fp, "/")
		fname := fArr[len(fArr)-1]
		fmt.Println("uploading", fname+"...")

		bucket.uploadFile(context.TODO(), fp, fname)

	}
}

func buildProject(rootDir string) {
	fmt.Println("building the project...")

	// npm install
	iCmd := exec.Command("npm", "install")
	iCmd.Dir = rootDir + "/output"
	iCmd.Stderr = os.Stderr
	iCmd.Stdout = os.Stdout

	runICmdErr := iCmd.Run()
	checkErr(runICmdErr)

	// npm run build
	bCmd := exec.Command("npm", "run", "build")
	bCmd.Dir = rootDir + "/output"
	bCmd.Stderr = os.Stderr
	bCmd.Stdout = os.Stdout

	runBCmdErr := bCmd.Run()
	checkErr(runBCmdErr)

	fmt.Println("build successful")
}
