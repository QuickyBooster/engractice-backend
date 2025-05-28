package models

import (
	"time"
)

type Test struct {
	Words       []Vocabulary       `json:"words"`
	Date        time.Time          `json:"date"`
	Correct     []Vocabulary       `json:"correct"`
	Wrong       []Vocabulary       `json:"wrong"`
	NearestMode bool               `json:"nearestMode"`
}
