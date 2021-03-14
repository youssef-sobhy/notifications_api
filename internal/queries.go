package internal

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx, database = connectToDb()
)

func connectToDb() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://notifications_mongo/"))
	if err != nil {
		// Report to rollbar or another logging service.
		fmt.Println(err.Error())
	}
	defer client.Disconnect(ctx)

	return ctx, client.Database("notifications_db")
}

func FetchNotifications() ([]Notification, err) {
	notificationsCollection := database.Collection("notifications")
	cursor, err := notificationsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var notifications []bson.M
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

func CreateNotifications(data []NotificationParams) ([]string, error) {
	var normalized_data []interface{}

	for _, obj := range data {
		for _, user := range obj.Users {
			obj.Notification.CreatedAt = time.Now()
			obj.Notification.UserID = user.ID
			obj.Notification.Platform = user.Platform
			normalized_data = append(normalized_data, obj.Notification)
		}
		Publish(obj.Notification, obj.Users)
	}

	res, err := notificationsCollection.InsertMany(ctx, normalized_data)
	if err != nil {
		return res.InsertedIDs, nil
	}
	return nil, err
}
