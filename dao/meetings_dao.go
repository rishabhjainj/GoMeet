package dao

import (
	"fmt"
	"log"
	"time"

	. "../models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MeetingsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "meetings"
)

func (m *MeetingsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

//find meetings list

func (m *MeetingsDAO) FindAll() ([]Meeting, error) {
	var meetings []Meeting
	err := db.C(COLLECTION).Find(bson.M{}).All(&meetings)
	return meetings, err
}

func (m *MeetingsDAO) FindByParticipant(participant string) ([]Meeting, error) {
	var meetings []Meeting
	err := db.C(COLLECTION).Find(bson.M{"participants._email": participant}).All(&meetings)
	return meetings, err
}

func (m *MeetingsDAO) MeetWithRSVPYes(participant string) ([]Meeting, error) {
	var meetings []Meeting
	err := db.C(COLLECTION).Find(bson.M{"participants._email": participant, "rsvp": "yes"}).All(&meetings)
	return meetings, err
}

func (m *MeetingsDAO) CheckForOverlappingMeetings(participant string, start time.Time, end time.Time) ([]Meeting, error) {
	var meetings []Meeting

	err := db.C(COLLECTION).Find(bson.M{
		"$or": []bson.M{
			bson.M{
				"$and": []bson.M{
					bson.M{"start_time": bson.M{
						"$lte": start,
					},
					},
					bson.M{"end_time": bson.M{
						"$gte": start,
					},
					},
				},
			},
			bson.M{
				"$and": []bson.M{
					bson.M{"start_time": bson.M{
						"$lte": end,
					},
					},
					bson.M{"end_time": bson.M{
						"$gte": end,
					},
					},
				},
			},
		},
		"participants._email": participant,
		"participants.rsvp":   "yes",
	}).All(&meetings)

	return meetings, err

}
func (m *MeetingsDAO) GetOverlappingMeetings(meetings []Meeting, start string, end string) {
	// for every meeting check for overlapping time

}
func (m *MeetingsDAO) FindByStartAndEnd(start string, end string) ([]Meeting, error) {
	layout := "2006-01-02T15:04:05"

	t_start, err := time.Parse(layout, start)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t_start)
	t_end, eerr := time.Parse(layout, end)
	if eerr != nil {
		fmt.Println(eerr)
	}
	fmt.Println(t_end.Hour())
	var meetings []Meeting

	nerr := db.C(COLLECTION).Find(bson.M{
		"start_time": bson.M{
			"$gte": t_start,
		},

		"end_time": bson.M{
			"$lte": t_end,
		},
	}).All(&meetings)
	return meetings, nerr
}

func (m *MeetingsDAO) Insert(meeting Meeting) error {
	err := db.C(COLLECTION).Insert(&meeting)
	return err
}

func (m *MeetingsDAO) FindById(id string) (Meeting, error) {
	var meeting Meeting
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&meeting)
	return meeting, err
}
