package models

type MgoXY struct {
	X float64 `json:"x" bson:"x" valid:"required"`
	Y float64 `json:"y" bson:"y" valid:"required"`
}
