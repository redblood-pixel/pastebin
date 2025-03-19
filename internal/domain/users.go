package domain

import "time"

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	Name      string `json:"name"`
	CreatedAt time.Time
	LastLogin time.Time
}
