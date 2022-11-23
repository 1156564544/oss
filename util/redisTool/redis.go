package redisTool

import (
	"context"
	"log"
	"os"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/go-redis/redis/v8"
)

const (
	// token的过期时间
	setDuration = 10
	// list的过期时间
	listDuration = 2
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
	rdb.Expire(context.Background(), list, listDuration*time.Second)
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

// 向Set中添加元素
func SetAdd(key string) {
	c, err := redigo.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Println("conn redis failed,", err)
		return
	}
	defer c.Close()
	_, err = c.Do("SET", key, 0, "EX", setDuration)
	if err != nil {
		log.Println(err)
		return
	}
}

// 查找集合中元素是否存在
func SetExist(key string) (bool, error) {
	c, err := redigo.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Println("conn redis failed,", err)
		return false, err
	}
	defer c.Close()
	return redigo.Bool(c.Do("EXISTS", key))
}

// 添加k-v元素到集合中
func AddKeyValue(key string, value string) error {
	c, err := redigo.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Println("conn redis failed,", err)
		return err
	}
	defer c.Close()
	_, err = c.Do("SET", key, value)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 从集合中获取k-v元素
func GetKeyValue(key string) (string, error) {
	c, err := redigo.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Println("conn redis failed,", err)
		return "", err
	}
	defer c.Close()
	return redigo.String(c.Do("GET", key))
}

// 删除集合中的元素
func DelKeyValue(key string) {
	c, err := redigo.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Println("conn redis failed,", err)
		return
	}
	defer c.Close()
	_, err = c.Do("DEL", key)
	if err != nil {
		log.Println(err)
		return
	}
}
