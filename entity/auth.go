package entity

// Auth ユーザー認証情報
type Auth struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}
