package entity

type Forum struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Author  string `json:"user"`
	Threads int64  `json:"threads,omitempty"`
	Posts   int64  `json:"posts,omitempty"`
}
