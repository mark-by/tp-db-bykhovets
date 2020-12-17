package entity

type User struct {
	Email    string `json:"email"`
	NickName string `json:"nickname"`
	FullName string `json:"fullname"`
	About    string `json:"about,omitempty"`
}

//easyjson:json
type UserList []User
