package models

type Vocabulary struct {
	ID         string   `json:"id" `
	English    string   `json:"english" `
	Vietnamese string   `json:"vietnamese" `
	Tag        []string `json:"tag" `
	Mp3        string   `json:"mp3" `
	CreatedAt  string   `json:"created_at" `
}
