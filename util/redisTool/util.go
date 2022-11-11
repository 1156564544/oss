package redisTool

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

// redis客户端
func redisConnect() (rdb *redis.Client) {

	var (
		redisServer string
		port        string
		password    string
	)

	redisServer = os.Getenv("RedisUrl")
	port = os.Getenv("RedisPort")
	password = os.Getenv("RedisPass")

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisServer + ":" + port,
		Password: password,
		DB:       0, // use default DB
	})

	return
}

// 广播消息
func PubMessage(channel, msg string) {
	rdb := redisConnect()
	rdb.Publish(context.Background(), channel, msg)
	rdb.Close()
}

// 订阅消息
func SubMessage(channel string) chan string {
	rdb := redisConnect()
	pubsub := rdb.Subscribe(context.Background(), channel)
	_, err := pubsub.Receive(context.Background())
	if err != nil {
		panic(err)
	}
	msg := make(chan string)
	ch := pubsub.Channel()

	go func() {
		for m := range ch {
			//fmt.Println(m.Channel, m.Payload)
			msg <- m.Payload
		}
	}()
	return msg
}

// 点对点发送消息
func PushMessage(list string, msg string) {
	rdb := redisConnect()
	rdb.LPush(context.Background(), list, msg)
	// fmt.Println(list, msg)
	rdb.Expire(context.Background(), list, 2*time.Second)
	rdb.Close()
}

// 点对点接收消息
func PopMessage(list string) string {
	rdb := redisConnect()
	defer rdb.Close()
	ip, err := rdb.LIndex(context.Background(), list, 0).Result()

	if err != nil {
		return ""
	}
	rdb.LPop(context.Background(), list)
	return ip
}
