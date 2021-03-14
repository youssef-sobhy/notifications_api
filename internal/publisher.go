package internal

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	},
}

var enqueuer = work.NewEnqueuer("notifications", redisPool)

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
			enqueuer.Enqueue(notification.Type, work.Q{"id": user.ID, "body": notification.Body})
		}
	}
}
