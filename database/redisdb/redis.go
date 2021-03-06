package redisdb

import (
	"GoLab/guard"
	"GoLab/server"
	"GoLab/tool"

	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	Client        *redis.Client
	RedisdbInfo   infoStruct
	valueFrom     string
	IsRedisEnable bool

	ctx = context.Background()
)

type infoStruct struct {
	URL      string
	Password string
}

func Set() {

	logString := "  Redis Info." + "\n"

	if server.Location == server.Cloud {
		if server.IsEnsaasServiceEnable && len(server.EnsaasService.Get("redis").MustArray()) != 0 {
			valueFrom = "ENSAAS_SERVICE"
			redisdb := server.EnsaasService.Get("redis").GetIndex(0).Get("credentials").MustMap()
			RedisdbInfo.URL = redisdb["host"].(string) + ":" + redisdb["port"].(json.Number).String()
			RedisdbInfo.Password = redisdb["password"].(string)
		} else {
			valueFrom = "ENV"
			RedisdbInfo.URL = os.Getenv("REDIS_URL")
			RedisdbInfo.Password = os.Getenv("REDIS_PASSWORD")
		}
	} else {
		valueFrom = "ENV"
		RedisdbInfo.URL = os.Getenv("REDIS_URL")
		redisdbPasswordFile := os.Getenv("REDIS_PASSWORD_FILE")
		if !tool.IsEmptyString(redisdbPasswordFile) {
			redisPassword, err := ioutil.ReadFile(redisdbPasswordFile)
			if err != nil {
				guard.Logger.Sugar().Fatalw("redis password file", "file path", redisdbPasswordFile)
			} else {
				RedisdbInfo.Password = string(redisPassword)
			}
		} else {
			RedisdbInfo.Password = os.Getenv("REDIS_PASSWORD")
		}
	}

	logString += "    FROM: " + valueFrom + "\n" +
		"      URL: " + RedisdbInfo.URL + "\n" +
		"      PASSWORD: " + RedisdbInfo.Password + "\n"

	fmt.Print(logString + "\n")

}

func Connect() {

	if !tool.IsEmptyString(RedisdbInfo.URL) {
		Client = redis.NewClient(&redis.Options{
			Addr:     RedisdbInfo.URL,
			Password: RedisdbInfo.Password,
		})

		_, err := Client.Ping(ctx).Result()
		if err != nil {
			IsRedisEnable = false
			guard.Logger.Fatal("redis login fail -> " + err.Error())
		} else {
			IsRedisEnable = true
		}
	} else {
		IsRedisEnable = false
	}

}
