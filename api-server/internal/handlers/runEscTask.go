package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/0xMishra/relay/api-server/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/joho/godotenv"
)

// loading the .env file
var err = godotenv.Load()

// all the necessary Environment variables
var (
	// these env used to configure ecs runtask command
	taskDefArn  = os.Getenv("TASK_DEF_ARN")
	clusterArn  = os.Getenv("CLUSTER_ARN")
	imageName   = os.Getenv("TASK_IMAGE_NAME")
	securityGrp = os.Getenv("SECUIRITY_GROUP")
	subnet1     = os.Getenv("SUBNET1")
	subnet2     = os.Getenv("SUBNET2")
	subnet3     = os.Getenv("SUBNET3")

	// these envs will be passed down to builder server
	bucketName   = os.Getenv("BUCKET_NAME")
	awsAccKeyId  = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecAccKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

	apiServerUrl = os.Getenv("API_SERVER_URL")
)

type Gitinfo struct {
	Url string `json:"url"`
	Pid string `json:"pid"`
}

type ResBody struct {
	Status string `json:"status"`
	Url    string `json:"url"`
}

func RunEcsTaskHandler(w http.ResponseWriter, r *http.Request) {
	var gitInfo Gitinfo
	err := json.NewDecoder(r.Body).Decode(&gitInfo)
	utils.CheckErr(err, true)

	defer r.Body.Close()

	PId := gitInfo.Pid

	runEcsTask(gitInfo, PId)

	res := ResBody{
		Status: "queued",
		Url:    PId + "." + apiServerUrl,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(res)
	utils.CheckErr(err, true)
}

func runEcsTask(gitInfo Gitinfo, pId string) {
	// AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	utils.CheckErr(err, true)

	client := ecs.NewFromConfig(cfg)

	// ecs run task config
	_, err = client.RunTask(context.TODO(), &ecs.RunTaskInput{
		Cluster:        aws.String(clusterArn),
		TaskDefinition: aws.String(taskDefArn),
		LaunchType:     types.LaunchType("FARGATE"),
		Count:          aws.Int32(1),
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				AssignPublicIp: types.AssignPublicIpEnabled,
				SecurityGroups: aws.ToStringSlice(
					[]*string{aws.String(securityGrp)},
				),
				Subnets: aws.ToStringSlice(
					[]*string{
						aws.String(subnet1),
						aws.String(subnet2),
						aws.String(subnet3),
					},
				),
			},
		},
		Overrides: &types.TaskOverride{
			ContainerOverrides: []types.ContainerOverride{
				{
					Name: aws.String(imageName),
					// envs for builder-server
					Environment: []types.KeyValuePair{
						{Name: aws.String("GIT_REPO_URL"), Value: aws.String(gitInfo.Url)},
						{Name: aws.String("PROJECT_ID"), Value: aws.String(pId)},
						{Name: aws.String("BUCKET_NAME"), Value: aws.String(bucketName)},
						{Name: aws.String("AWS_ACCESS_KEY_ID"), Value: aws.String(awsAccKeyId)},
						{
							Name:  aws.String("AWS_SECRET_ACCESS_KEY"),
							Value: aws.String(awsSecAccKey),
						},
						{Name: aws.String("REDIS_URL"), Value: aws.String(utils.RedisUrl)},
					},
				},
			},
		},
	})

	utils.CheckErr(err, true)
}
