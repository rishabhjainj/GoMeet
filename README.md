## Go Meet

Meeting Scheduling API STACK

### Technology Stack

> MongoDB  
> GoLang

### Setup Instructions

1.  Clone the repository
2.  Install following packages  
     `go get gopkg.in/mgo.v2 `
3.  Setup MongoDB
    Server : localhost
    Database : meetings_db
    Collection : meetings
    Port currently set to default mongoDB port : 27017
4.  To run
    `go build`
    `go run app.go`
5.  Available Endpoints
    a. GET : http://localhost:27017/meetings/ : Returns List of all meetings
    Example :
    [
    {
    "id": "5f4eab2b62b33740b8dd0a85",
    "title": "6to8 meeting",
    "created_at": "2020-09-02T01:42:27.62+00:00",
    "start_time": "2020-09-01T22:34:05+00:00",
    "end_time": "2020-09-01T22:34:05+00:00",
    "participants": [
    {
    "name": "rishabh",
    "email": "jrishabh252",
    "rsvp": "yes"
    }
    ]
    },
    ]

    b. POST : http://localhost:27017/meetings/ : Creates a new meeting after checking if there are any overlapping
    meetings in which case it will show an error

    Example :
    BODY:

            {

            "title": "6to8 meeting",
            "created_at": "2020-09-01T22:34:05+00:00",
            "start_time": "2020-09-01T22:34:05+00:00",
            "end_time": "2020-09-01T22:34:05+00:00",
            "participants": [
                  {
                  "name": "rishabh",
                  "email": "jrishabh252",
                  "rsvp": "yes"
                  }
            ]
            }

    IMPORTANT !!! Keep the time format as : "2020-09-01T22:34:05+00:00"

    c. GET : http://localhost:27017/meetings/?participant="YOUR_PARTICIPANT_EMAIL" : Returns list of meetings of that particular participant

    Example : http://localhost:27017/meetings/?participant=jrishabh252

    Result :
    [
    {
    "id": "5f4eae3b62b33750bcd2d7f9",
    "title": "6to8 meeting",
    "created_at": "2020-09-02T01:55:31.69+05:30",
    "start_time": "2020-09-02T04:04:05+05:30",
    "end_time": "2020-09-02T04:04:05+05:30",
    "participants": [
    {
    "name": "rishabh",
    "email": "jrishabh252",
    "rsvp": "yes"
    }
    ]
    },
    ]
    d. GET : http://localhost:27017/meetings/?start=YOUR_START_TIME"&end=YOUR_END_TIME : Returns list of meetings within start and end timings (inclusive of these two time points)

    Example : http://localhost:27017/meetings/?start=2020-09-01T22:34:05&end=2020-09-01T22:34:05

    !!!IMPORTANT : Keep the time format in url as : 2020-09-01T22:34:05

    Result :

    [
    {
    "id": "5f4eae3b62b33750bcd2d7f9",
    "title": "6to8 meeting",
    "created_at": "2020-09-02T01:55:31.69+05:30",
    "start_time": "2020-09-02T04:04:05+05:30",
    "end_time": "2020-09-02T04:04:05+05:30",
    "participants": [
    {
    "name": "rishabh",
    "email": "jrishabh252",
    "rsvp": "yes"
    }
    ]
    },
    ]
