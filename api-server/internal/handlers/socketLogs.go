package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/0xMishra/relay/api-server/internal/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// this upgrades http header to use websocket protocol
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// open for anyone (CORS policy)
		return true
	},
}

func SocketLogsHandler(w http.ResponseWriter, r *http.Request) {
	rdb, ok := r.Context().Value("redisClient").(*redis.Client)
	if !ok {
		utils.CheckErr(errors.New("can't initialize redis client"), true)
	}

	PId := strings.Split(r.URL.Path, "/")[2]
	// establishing a socket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.CheckErr(err, false)

	fmt.Println("user connected")

	defer conn.Close()

	// susbscring to redis channel to get realtime log from builder-server with redis
	ctx := context.Background()
	sub := rdb.Subscribe(ctx, "log:"+PId)
	ch := sub.Channel()

	fmt.Println("Subscribed to channel: "+"log:"+PId+",", sub)

	defer sub.Close()

	for msg := range ch {
		fmt.Println(string(msg.Payload))

		err = conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		if err != nil {
			fmt.Println("error in socket connection while sending messages", err)
		}
	}
}
