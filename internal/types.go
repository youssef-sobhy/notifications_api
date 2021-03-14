package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Body       string             `bson:"body,required,omitempty"`
	CreatedAt  time.Time          `bson:"created_at,omitempty"`
	Type       string             `bson:"type,required,omitempty"`
	Platform   int                `bson:"platform,omitempty"`
	Title      string             `bson:"title,omitempty"`
	UserID     string             `bson:"user_id,required,omitempty"`
	ExternalID int                `bson:"external_id,required,omitempty"`
}

type User struct {
	ID         string `json:"id,required"`
	Platform   int    `json:"platform,required"`
	ExternalID int    `json:"external_id,required"`
}

type UserQueryParams struct {
	ExternalID int `uri:"external_id"`
}

type SendParams struct {
	Data []NotificationParams `json:"data,required"`
}

type NotificationParams struct {
	Notification Notification `json:"notification,required"`
	Users        []User       `json:"users,required"`
}
