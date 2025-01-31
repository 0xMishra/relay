package middlewares

import (
	"context"
	"net/http"

	"github.com/0xMishra/relay/api-server/internal/utils"
	"github.com/go-redis/redis/v8"
)

func RedisSetup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr, err := redis.ParseURL(utils.RedisUrl)
		utils.CheckErr(err, false)

		rdb := redis.NewClient(addr)

		err = rdb.Ping(context.Background()).Err()
		utils.CheckErr(err, false)

		ctx := context.WithValue(r.Context(), "redisClient", rdb)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
