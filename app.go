package main

import (
	. "./dao"
	. "./models"
	// myRouter "./router"
	"encoding/json"
	"regexp"

	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strings"
	"time"
)

var dao = MeetingsDAO{}

func getMeetingsByParticipant(w http.ResponseWriter, participant string) {
	meetings, err := dao.FindByParticipant(participant)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, meetings)
	return
}

func getMeetingsByTimeRange(w http.ResponseWriter, start string, end string) {
	meetings, err := dao.FindByStartAndEnd(start, end)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, meetings)
	return
}

//router for get requests

func getAllMeetings(w http.ResponseWriter) {
	meetings, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, meetings)
}

func CreateMeeting(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var meeting Meeting
	if err := json.NewDecoder(r.Body).Decode(&meeting); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}
	//before inserting check for overlapping meetings

	// do this for every participant if their rsvp is yes
	for _, v := range meeting.Participants {
		if v.RSVP == "yes" {
			participantMeetings, err := dao.CheckForOverlappingMeetings(v.Email, meeting.StartTime, meeting.EndTime)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if len(participantMeetings) > 0 {
				respondWithError(w, http.StatusInternalServerError, "There are already scheduled meetings")
				return
			}
		}

	}

	meeting.ID = bson.NewObjectId()
	meeting.CreatedAt = time.Now()
	if err := dao.Insert(meeting); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, meeting)

}

func GetMeetingByID(w http.ResponseWriter, id string) {

	if id != "" {
		meeting, err := dao.FindById(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid Meeting ID")
			return
		}
		respondWithJson(w, http.StatusOK, meeting)
		return
	}
}
func GetMeetingsOfParticipant(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func init() {
	// config.Read()
	dao.Server = "localhost"
	dao.Database = "meetings_db"
	dao.Connect()
}

func multiplexer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getMux(w, r)
	case "POST":
		postMux(w, r)
	}
}
func postMux(w http.ResponseWriter, r *http.Request) {
	CreateMeeting(w, r)
}
func getMux(w http.ResponseWriter, r *http.Request) {
	//check if request contains id
	var checkId = regexp.MustCompile(`^\d*[a-zA-Z][a-zA-Z\d]*$`)
	str1 := strings.Split(r.URL.Path, "/")
	id := str1[len(str1)-1:][0]
	if checkId.MatchString(id) {
		GetMeetingByID(w, id)
	}
	//if not then check if participant is present

	v := r.URL.Query()
	participant := v.Get("participant")
	if participant != "" {
		getMeetingsByParticipant(w, participant)
		return
	}
	start := v.Get("start")
	end := v.Get("end")
	if start != "" {
		if end != "" {
			//if both start and end times are present then search
			getMeetingsByTimeRange(w, start, end)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "end time is required")
		return
	}
	if end != "" {
		if start != "" {
			getMeetingsByTimeRange(w, start, end)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "start time is required")
		return
	}
	//finally throw all meetings
	getAllMeetings(w)

}

func main() {
	// r := myRouter.Newrouter()
	// r.GET("/meetings", GetAllMeetingsEndPoint)
	// r.POST("/meetings", CreateMeeting)
	// if err := http.ListenAndServe(":27017", r); err != nil {
	// 	log.Fatal(err)
	// }

	http.HandleFunc("/", multiplexer)
	if err := http.ListenAndServe(":27017", nil); err != nil {
		log.Fatal(err)
	}

}
