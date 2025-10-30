package redis

import (
	"context"

	"github.com/genda/genda-api/pkg/config"
	"github.com/go-redis/redis/v8"

	"log"
	"os"
	"strconv"
)

type Redis *redis.Client

func NewConnection() (*redis.Client, error) {
	log := log.New(os.Stdout, "[internal.storage.redis] ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)
	conf := config.New()

	redisDb, _ := strconv.Atoi(conf.RedisDB.Db)
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.RedisDB.Host,
		Password: conf.RedisDB.Password,
		DB:       redisDb,
	})

	var ctx = context.Background()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatal("Error occurred when connecting to redis: ", err)
		return nil, err
	} else {
		log.Println("Connected to redis")
	}

	return rdb, nil
}
