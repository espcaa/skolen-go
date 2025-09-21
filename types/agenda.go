package types

import "time"

type TimetableDay struct {
	Date        time.Time    `json:"date"`
	Lessons     []Lesson     `json:"lessons"`
	Assignments []Assignment `json:"assignments"`
}

type Lesson struct {
	ID                          string    `json:"id"`
	StartDateTime               time.Time `json:"startDateTime"`
	EndDateTime                 time.Time `json:"endDateTime"`
	Location                    string    `json:"location"`
	Canceled                    bool      `json:"canceled"`
	AnyContent                  bool      `json:"anyContent"`
	AnyHomeworkToDoForTheLesson bool      `json:"anyHomeworkToDoForTheLesson"`
	AnyHomeworkToDoAfterLesson  bool      `json:"anyHomeworkToDoAfterLesson"`
	Subject                     Subject   `json:"subject"`
	Teachers                    []Teacher `json:"teachers"`
}

type Subject struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Color string `json:"color"`
}

type Teacher struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoUrl"`
}
