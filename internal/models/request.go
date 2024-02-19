package models

type Request struct {
	ID         uint     `json:"id"`
	Method     string   `json:"method"`
	Path       string   `json:"path"`
	GetParams  []Params `json:"get_params"`
	PostParams []Params `json:"post_params"`
	Headers    []Params `json:"headers"`
	Cookie     []Params `json:"cookie"`
	Body       string   `json:"body"`
}

type Params struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
