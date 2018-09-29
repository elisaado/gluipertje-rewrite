package models

// User struct
type User struct {
	ID int `json:"id" storm:"id,increment"`
	Model
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}
