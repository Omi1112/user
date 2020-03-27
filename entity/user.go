package entity

// User ユーザー情報
type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"    binding:"required,email"`
	Name     string `json:"name"     binding:"required,max=100"`
	Password string `json:"password" binding:"required,min=8,max=100,alphanum"`
}
