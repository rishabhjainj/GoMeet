package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Meeting struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	Title        string        `bson:"title" json:"title"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	StartTime    time.Time     `bson:"start_time" json:"start_time"`
	EndTime      time.Time     `bson:"end_time" json:"end_time"`
	Participants []Participant `bson:"participants" json:"participants"`
}

type Participant struct {
	Name  string `bson:"name" json:"name"`
	Email string `bson:"_email" json:"email"`
	RSVP  string `bson:"rsvp" json:"rsvp"`
}

func (u *Meeting) GetBSON() (interface{}, error) {
	time.Local = time.UTC
	u.CreatedAt = time.Now()
	type my *Meeting
	return my(u), nil
}
