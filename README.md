## Demo
https://ts9fugmle2.ufs.sh/f/334k6TPdXjoiOs8CTlRJiX2kPVs6rdat8CbTO7EunjDIZqcW

## Architecture

![diagram-export-1-31-2025-11_35_26-PM](https://github.com/user-attachments/assets/584b3c8d-e66f-4904-b8f2-39b05acfb4b3)

## How to run the project

- Build the build-server docker image and push it to your AWS ECR
- Go to your AWS console and create a new ECS cluster
- Then create a task definition in that cluster using the build-server docker image from your AWS ECR
- Now create a `.env` file inside each top-level directory and add all the required environment variables ( take a look at the `.env.example` for reference)
- After that `cd` into api-server and run `go run cmd/server/main.go`
- Now run the client project with `npm i` then `npm run dev`

## PORTS

- api-server:`3000`
- client:`3001`
