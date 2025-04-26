package models

import (
	"time"
)

type Test struct {
	ID      string       `json:"id" bson:"_id"`
	Words   []Vocabulary `json:"words"`
	Date    time.Time    `json:"date"`
	Correct []Vocabulary `json:"correct"`
	Wrong   []Vocabulary `json:"wrong"`
}
