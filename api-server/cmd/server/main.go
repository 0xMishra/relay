package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Pallinder/go-randomdata"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// loading .env file
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
	redisUrl     = os.Getenv("REDIS_URL")
)

func ecsConfig(gitInfo Gitinfo, pId string) {
	// AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	checkErr(err, true)

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
						{Name: aws.String("REDIS_URl"), Value: aws.String(redisUrl)},
					},
				},
			},
		},
	})

	checkErr(err, true)
}

func checkErr(err error, fatal bool) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if fatal {
			os.Exit(1)
		}
	}
}

type Gitinfo struct {
	Url string `json:"url"`
	Pid string `json:"pid"`
}

type ResBody struct {
	Status string `json:"status"`
	Url    string `json:"url"`
}

var (
	rdb *redis.Client
	ctx = context.Background()
	pId = randomdata.Adjective()
)

// this upgrades http header to use websocket protocol
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// open for anyone (CORS policy)
		return true
	},
}

func main() {
	checkErr(err, false)
	rdb = redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	mux := http.NewServeMux()

	mux.HandleFunc("POST /project", runEcsTaskHandler)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// establishing a socket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		checkErr(err, false)

		fmt.Println("user connected")

		defer conn.Close()

		// susbscring to redis channel to get realtime log from builder-server with redis
		ctx = context.Background()
		sub := rdb.Subscribe(ctx, "log:"+pId)
		ch := sub.Channel()

		fmt.Println(pId, sub, ch)

		defer sub.Close()

		for msg := range ch {
			fmt.Println(string(msg.Payload))

			err = conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				fmt.Println("error in socket connection while sending messages", err)
			}
		}
	})

	fmt.Println("api server running on PORT:3000")
	http.ListenAndServe("localhost:3000", mux)
}

func runEcsTaskHandler(w http.ResponseWriter, r *http.Request) {
	var gitInfo Gitinfo
	err := json.NewDecoder(r.Body).Decode(&gitInfo)
	checkErr(err, true)

	defer r.Body.Close()

	if len(gitInfo.Pid) > 0 {
		pId = gitInfo.Pid
	}

	ecsConfig(gitInfo, pId)

	res := ResBody{
		Status: "queued",
		Url:    "http://" + pId + ".localhost:8000",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(res)
	checkErr(err, true)
}
