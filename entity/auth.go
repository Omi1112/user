package entity

// Auth ユーザー認証情報を保管
type Auth struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}
