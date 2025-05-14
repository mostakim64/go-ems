package conn

import (
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/golang-course-utils/logger"

	"github.com/go-redis/redis"
)

var client *redis.Client

func ConnectRedis() {
	conf := config.Redis()

	logger.Info("connecting to redis at ", conf.Host, ":", conf.Port, "...")

	client = redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Password: conf.Pass,
		DB:       conf.Db,
	})

	if _, err := client.Ping().Result(); err != nil {
		logger.Error("failed to connect redis: ", err)
		panic(err)
	}

	logger.Info("redis connection successful...")
}

type RedisClient struct{}

func NewRedisClient() *RedisClient {
	return &RedisClient{}
}
