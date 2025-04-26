package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Test struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Words       []Vocabulary       `json:"words"`
	Date        time.Time          `json:"date"`
	Correct     []Vocabulary       `json:"correct"`
	Wrong       []Vocabulary       `json:"wrong"`
	NearestMode bool               `json:"nearestMode"`
}
