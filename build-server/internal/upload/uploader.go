package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-redis/redis/v8"
)

const BUCKET_URL = ""

var (
	rdb *redis.Client
	ctx = context.Background()
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		publishRedisLog(err.Error())
	}
}

func publishRedisLog(log string) {
	pid := os.Getenv("PROJECT_ID")
	ctx = context.Background()
	rdb.Publish(ctx, "log:"+pid, log)
}

type Bucket struct {
	S3Client *s3.Client
}

func (bucket Bucket) uploadFile(
	ctx context.Context,
	fPath string,
	fileName string,
	fileExt string,
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
	mimetype := mime.TypeByExtension(fileExt)

	bucketName := bName
	objectKey := "_output/" + pId + "/" + fileName

	largeBuffer := bytes.NewReader(data)
	var partMiBs int64 = 10

	// uploading file in small chunks
	uploader := manager.NewUploader(bucket.S3Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        largeBuffer,
		ContentType: &mimetype,
	})
	if err != nil {
		log.Printf("Error while uploading %s, %s\n", fileName, err)
		publishRedisLog(fmt.Sprintf("Error while uploading %s, %s\n", fileName, err))
		return err
	}

	fmt.Printf(
		"Successfully uploaded %s\n",
		fileName,
	)
	publishRedisLog(
		"Successfully uploaded " + fileName,
	)

	return nil
}

func Init() {
	redisUrl := os.Getenv("REDIS_URL")
	addr, err := redis.ParseURL(redisUrl)
	checkErr(err)

	rdb = redis.NewClient(addr)

	err = rdb.Ping(ctx).Err()
	checkErr(err)

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
				publishRedisLog(err.Error())
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

	fmt.Println("\ndeploying project...")
	publishRedisLog("\ndeploying project...")

	// uploading files one by one
	for _, fp := range filePaths {
		fArr := strings.Split(fp, "/")
		fname := fArr[len(fArr)-1]

		fileExt := filepath.Ext(fp)

		if strings.Contains(fp, "assets") {
			fname = "assets/" + fname
		}

		fmt.Println("\nuploading ", fname+"...")
		publishRedisLog("\nuploading " + fname + "...")

		bucket.uploadFile(context.TODO(), fp, fname, fileExt)
	}

	fmt.Printf("deployment successful")
	publishRedisLog("deployment successful")
}

func buildProject(rootDir string) {
	fmt.Println("\ninstalling packages...")
	publishRedisLog("\ninstalling packages...")

	// npm install
	iCmd := exec.Command("npm", "install")
	iCmd.Dir = rootDir + "/output"
	iCmd.Stderr = os.Stderr
	iCmd.Stdout = os.Stdout

	runICmdErr := iCmd.Run()
	checkErr(runICmdErr)

	fmt.Println("\nbuilding the project...")
	publishRedisLog("\nbuilding the project...")

	// npm run build
	bCmd := exec.Command("npm", "run", "build")
	bCmd.Dir = rootDir + "/output"
	bCmd.Stderr = os.Stderr
	bCmd.Stdout = os.Stdout

	runBCmdErr := bCmd.Run()
	checkErr(runBCmdErr)

	fmt.Println("build successful")
	publishRedisLog("build successful")
}
