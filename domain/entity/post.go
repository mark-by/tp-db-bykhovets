package entity

import "github.com/jackc/pgtype"

type Post struct {
	ID       int64            `json:"-"`
	Message  string           `json:"message"`
	IsEdited bool             `json:"isEdited"`
	Parent   int64            `json:"parent"`
	Created  string           `json:"created"`
	Author   string           `json:"author"`
	Thread   int32            `json:"thread"`
	Forum    string           `json:"forum"`
	Path     pgtype.Int8Array `json:"-"`
}

type PostFull struct {
	Post   Post   `json:"post"`
	Author User   `json:"author"`
	Thread Thread `json:"thread"`
	Forum  Forum  `json:"forum"`
}
