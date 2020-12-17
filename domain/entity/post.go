package entity

import "github.com/jackc/pgtype"

type Post struct {
	ID       int64            `json:"id"`
	Message  string           `json:"message"`
	IsEdited bool             `json:"isEdited,omitempty"`
	Parent   int64            `json:"parent"`
	Created  string           `json:"created"`
	Author   string           `json:"author"`
	Thread   int32            `json:"thread"`
	Forum    string           `json:"forum"`
	Path     pgtype.Int8Array `json:"-"`
}

//easyjson:json
type PostList []Post

type PostFull struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}
