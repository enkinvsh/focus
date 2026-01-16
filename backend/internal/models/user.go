package models

import "time"

type User struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username,omitempty"`
	FirstName  string    `json:"first_name,omitempty"`
	Language   string    `json:"language"`
	Timezone   string    `json:"timezone"`
	ThemeIndex int       `json:"theme_index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
