package models

import "time"

type Vocabulary struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Owner     string    `json:"owner,omitempty" bson:"owner,omitempty"`
	Japanese  Japanese  `json:"japanese"`
	Thai      []string  `json:"thai"`
	English   []string  `json:"english"`
	Examples  []string  `json:"examples"`
	Image     string    `json:"string"`
	Voice     string    `json:"vocie"`
	Type      string    `json:"type"`
	Tags      []string  `json:"tags"`
	IsShow    bool      `json:"isShow"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Japanese struct {
	Kanji  string `json:"kanji"`
	Kana   string `json:"kana"`
	Romaji string `json:"romaji"`
}
