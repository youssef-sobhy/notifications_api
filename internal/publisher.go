package internal

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var (
	redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}
	rdb = redisv8.NewClient(&redisv8.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	ctx      = context.Background()
	enqueuer = work.NewEnqueuer("notifications", redisPool)
	smsLimit = os.Getenv("SMS_REQUEST_LIMIT")
)

type Sent struct {
	Time  time.Time
	Count int
}

func Publish(notification Notification, users []User) {
	if notification.Type == "push" {
		platforms := map[int][]string{}
		for _, user := range users {
			if _, ok := platforms[user.Platform]; ok == false {
				platforms[user.Platform] = []string{}
			}
			platforms[user.Platform] = append(platforms[user.Platform], user.ID)
		}
		for platform, tokens := range platforms {
			enqueuer.Enqueue(notification.Type, work.Q{"tokens": tokens, "body": notification.Body,
				"title": notification.Title, "platform": platform})
		}
	} else {
		for _, user := range users {
			publishSms(notification, user)
		}
	}
}

func publishSms(notification Notification, user User) {
	var sent Sent
	var secondsFromNow int64
	var ttl time.Duration
	mu := sync.Mutex{}
	mu.Lock()
	val, err := rdb.Get(ctx, "sent").Result()
	json.Unmarshal([]byte(val), sent)
	mu.Unlock()

	if err != nil {
		sent.Time = time.Now()
		sent.Count = 1
	} else {
		limit, _ := strconv.Atoi(smsLimit)
		if sent.Count >= limit {
			secondsFromNow = int64(60 - sent.Time.Second())
			sent.Time = sent.Time.Add(time.Minute * 1)
			sent.Count = 1
		} else {
			sent.Count += 1
			secondsFromNow = 0
		}
	}

	currentTime := time.Now()

	if sent.Time.Before(currentTime) {
		ttl = 60
	} else {
		ttl = time.Duration(sent.Time.Sub(currentTime).Seconds())
	}

	mu.Lock()
	data, _ := json.Marshal(sent)
	rdb.Set(ctx, "sent", string(data), ttl)
	mu.Unlock()
	enqueuer.EnqueueIn(notification.Type, secondsFromNow, work.Q{"id": user.ID, "body": notification.Body})
}
