package internal

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectToDb() (context.Context, *mongo.Client) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017/?connect=direct"))
	if err != nil {
		// Report to rollbar or another logging service.
		fmt.Println(err.Error())
	}

	return ctx, client
}

func FetchNotifications() ([]primitive.M, error) {
	ctx, client := connectToDb()
	database := client.Database("notifications_db")
	notificationsCollection := database.Collection("notifications")
	cursor, err := notificationsCollection.Find(ctx, bson.M{})
	defer client.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	var notifications []bson.M
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

func CreateNotifications(data []NotificationParams) ([]primitive.M, error) {
	ctx, client := connectToDb()
	database := client.Database("notifications_db")
	notificationsCollection := database.Collection("notifications")
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
	defer client.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	var notifications []bson.M
	query := bson.M{"_id": bson.M{"$in": res.InsertedIDs}}
	cursor, err := notificationsCollection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}
