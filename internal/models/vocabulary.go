package models

type Vocabulary struct {
	English    string `json:"english"`
	Vietnamese string `json:"vietnamese"`
	MP3        string `json:"mp3"`
	Tag        string `json:"tag"`
	Point      int    `json:"point"`
}
