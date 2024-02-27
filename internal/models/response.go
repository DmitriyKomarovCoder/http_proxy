package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Response struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code    int                `json:"code"`
	Headers []Params           `json:"headers"`
	Body    string             `json:"body"`
	Cookie  []Params           `json:"cookie"`
}
