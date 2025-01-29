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
)

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
		log.Printf("Error while uploading object to %s. %s\n", bucketName, err)
		return err
	}

	fmt.Printf(
		"Successfully uploaded file %s",
		fPath,
	)

	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
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

		fileExt := filepath.Ext(fp)

		if strings.Contains(fp, "assets") {
			fname = "assets/" + fname
		}

		fmt.Println("\n\nuploading", fname+"...")

		bucket.uploadFile(context.TODO(), fp, fname, fileExt)

	}
}

func updateBuildPath(filePath string, pid string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read index.html: %w", err)
	}

	// Replace base path to _output/{pId}
	updatedContent := strings.ReplaceAll(
		string(content),
		"vite build",
		"vite build --base=/_output/"+pid+"/",
	)

	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to update index.html: %w", err)
	}

	return nil
}

func buildProject(rootDir string) {
	fmt.Println("\ninstalling packages...")

	pId := os.Getenv("PROJECT_ID")
	/*
	   we are changing the vite build to vite build --base="/_output/{pId}/ since thats
	   our nested path of static files in S3 bucket
	*/
	updateBuildPath(rootDir+"/output/package.json", pId)

	// npm install
	iCmd := exec.Command("npm", "install")
	iCmd.Dir = rootDir + "/output"
	iCmd.Stderr = os.Stderr
	iCmd.Stdout = os.Stdout

	runICmdErr := iCmd.Run()
	checkErr(runICmdErr)

	fmt.Println("\nbuilding the project...")

	// npm run build
	bCmd := exec.Command("npm", "run", "build")
	bCmd.Dir = rootDir + "/output"
	bCmd.Stderr = os.Stderr
	bCmd.Stdout = os.Stdout

	runBCmdErr := bCmd.Run()
	checkErr(runBCmdErr)

	fmt.Println("build successful")
}
