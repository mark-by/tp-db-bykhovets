package entity

type User struct {
	ID       int32  `json:"-"`
	Email    string `json:"email"`
	NickName string `json:"nickname"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
}
