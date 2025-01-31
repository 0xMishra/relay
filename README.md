## Architecture
![arch](https://github.com/user-attachments/assets/0b80ebaf-f3aa-4055-a87b-d931288f5452)

## How to run the project
- Build the build-server docker image and push it to your AWS ECR
- Go to your AWS console and create a new ECS cluster
- Then create a task definition in that cluster using the build-server docker image from your AWS ECR
- Now create a `.env` file inside each top-level directory and add all the required environment variables ( take a look at the `.env.example` for reference)
- After that `cd` into api-server and reverse-proxy in a different terminal and run `go run cmd/server/main.go` in both directories
- Now run the client project with `npm i` then `npm run dev`

## PORTS
- api-server:`3000`
- reverse-proxy:`8000`
- client:`3001`
