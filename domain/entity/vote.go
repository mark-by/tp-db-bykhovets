package entity

type Vote struct {
	Voice  int32  `json:"voice"`
	Author string `json:"nickname"`
}
