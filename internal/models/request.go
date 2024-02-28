package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Request struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Scheme     string             `json:"scheme"`
	Method     string             `json:"method"`
	Path       string             `json:"path"`
	Host       string             `json:"host"`
	GetParams  []Params           `json:"get_params"`
	PostParams []Params           `json:"post_params"`
	Headers    []Params           `json:"headers"`
	Cookie     []Params           `json:"cookie"`
	Body       string             `json:"body"`
}

type Params struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
