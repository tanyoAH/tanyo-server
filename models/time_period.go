package models

import "time"

type TimePeriod struct {
	Start time.Time `json:"start" bson:"start" valid:"required"`
	End   time.Time `json:"end" bson:"end" valid:"required"`
}
