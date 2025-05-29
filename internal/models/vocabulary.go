package models

type Vocabulary struct {
	// Order to remember the position of the word in the list
	Order      int    `json:"order"`
	English    string `json:"english"`
	Vietnamese string `json:"vietnamese"`
	MP3        string `json:"mp3"`
	Tag        string `json:"tag"`
	Point      int    `json:"point"`
}
