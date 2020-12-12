package entity

type Vote struct {
	ID     int32  `json:"-"`
	Voice  int32  `json:"voice"`
	Author string `json:"nickname"`
	Thread int32  `json:"-"`
}
