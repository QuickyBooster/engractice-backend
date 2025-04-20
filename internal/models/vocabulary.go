package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vocabulary struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	English    string             `json:"english" `
	Vietnamese string             `json:"vietnamese" `
	Tag        []string           `json:"tag" `
	Mp3        string             `json:"mp3" `
	CreatedAt  time.Time          `json:"created_at" `
}
