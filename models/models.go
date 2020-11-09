package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Article ...
type Article struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title             string             `json:"title,omitempty" bson:"title,omitempty"`
	SubTitle          string             `json:"subtitle" bson:"subtitle,omitempty"`
	Content           string             `json:"content" bson:"content,omitempty"`
	CreationTimeStamp time.Time          `json:"creationtime" bson:"crationtime,omitempty"`
}
