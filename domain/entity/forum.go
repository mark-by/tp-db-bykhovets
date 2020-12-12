package entity

type Forum struct {
	ID      int32  `json:"-"`
	Slug    string `json:"slug"`
	Author  string `json:"user"`
	Title   string `json:"title"`
	Threads int64  `json:"threads"`
	Posts   int64  `json:"posts"`
}
