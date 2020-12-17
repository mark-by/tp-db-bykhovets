package entity

type Thread struct {
	ID      int32  `json:"id"`
	Slug    string `json:"slug,omitempty"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Created string `json:"created"`
	Votes   int32  `json:"votes"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
}
