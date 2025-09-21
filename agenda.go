package skolengo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/espcaa/skolen-go/types"
	utils "github.com/espcaa/skolen-go/utils"
)

type BaseResponse struct {
	Data     []map[string]interface{} `json:"data"`
	Included []map[string]interface{} `json:"included"`
}

func (c *Client) GetTimetable(userID, schoolID, emsCode string, periodStart, periodEnd time.Time, limit int) ([]types.TimetableDay, error) {
	if periodStart.IsZero() {
		periodStart = time.Now()
	}
	if periodEnd.IsZero() {
		periodEnd = periodStart.AddDate(0, 1, 0)
	}
	if limit == 0 {
		limit = 50
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/agendas", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("filter[student.id]", userID)
	q.Add("filter[date][GE]", periodStart.Format("2006-01-02"))
	q.Add("filter[date][LE]", periodEnd.Format("2006-01-02"))
	q.Add("include", "lessons,lessons.subject,lessons.teachers,homeworkAssignments,homeworkAssignments.subject")
	q.Add("page[limit]", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+c.TokenSet.AccessToken)
	req.Header.Set("x-skolengo-ems-code", emsCode)
	req.Header.Set("x-skolengo-school-id", schoolID)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var baseResp BaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&baseResp); err != nil {
		return nil, err
	}

	includedMap := make(map[string]map[string]interface{})
	for _, item := range baseResp.Included {
		key := fmt.Sprintf("%s:%s", item["type"], item["id"])
		includedMap[key] = item
	}

	var result []types.TimetableDay
	for _, day := range baseResp.Data {
		rawLessons := utils.GetMultipleRelations(day["relationships"].(map[string]interface{})["lessons"].(map[string]interface{}))
		rawAssignments := utils.GetMultipleRelations(day["relationships"].(map[string]interface{})["homeworkAssignments"].(map[string]interface{}))

		var lessons []types.Lesson
		var assignments []types.Assignment

		for _, rawLesson := range rawLessons {
			lessonID := rawLesson["id"].(string)
			lessonData := includedMap["lesson:"+lessonID]

			subjectID := utils.GetSingleRelation(lessonData["relationships"].(map[string]interface{})["subject"].(map[string]interface{}))["id"].(string)
			subjectData := includedMap["subject:"+subjectID]

			var teachers []types.Teacher
			teachersRel := utils.GetMultipleRelations(lessonData["relationships"].(map[string]interface{})["teachers"].(map[string]interface{}))
			for _, t := range teachersRel {
				tData := includedMap["teacher:"+t["id"].(string)]
				teachers = append(teachers, types.Teacher{
					ID:        t["id"].(string),
					Title:     tData["attributes"].(map[string]interface{})["title"].(string),
					FirstName: tData["attributes"].(map[string]interface{})["firstName"].(string),
					LastName:  tData["attributes"].(map[string]interface{})["lastName"].(string),
					PhotoURL:  tData["attributes"].(map[string]interface{})["photoUrl"].(string),
				})
			}

			lessons = append(lessons, types.Lesson{
				ID:                          lessonID,
				StartDateTime:               parseTime(lessonData["attributes"].(map[string]interface{})["startDateTime"].(string)),
				EndDateTime:                 parseTime(lessonData["attributes"].(map[string]interface{})["endDateTime"].(string)),
				Location:                    lessonData["attributes"].(map[string]interface{})["location"].(string),
				Canceled:                    lessonData["attributes"].(map[string]interface{})["canceled"].(bool),
				AnyContent:                  lessonData["attributes"].(map[string]interface{})["anyContent"].(bool),
				AnyHomeworkToDoForTheLesson: lessonData["attributes"].(map[string]interface{})["anyHomeworkToDoForTheLesson"].(bool),
				AnyHomeworkToDoAfterLesson:  lessonData["attributes"].(map[string]interface{})["anyHomeworkToDoAfterLesson"].(bool),
				Subject: types.Subject{
					ID:    subjectID,
					Label: subjectData["attributes"].(map[string]interface{})["label"].(string),
					Color: subjectData["attributes"].(map[string]interface{})["color"].(string),
				},
				Teachers: teachers,
			})
		}

		for _, rawAssignment := range rawAssignments {
			assignID := rawAssignment["id"].(string)
			assignData := includedMap["homework:"+assignID]

			subjectRel := utils.GetSingleRelation(assignData["relationships"].(map[string]interface{})["subject"].(map[string]interface{}))
			subjectID := ""
			subjectData := map[string]interface{}{}
			if subjectRel != nil {
				subjectID = subjectRel["id"].(string)
				subjectData = includedMap["subject:"+subjectID]
			}

			assignments = append(assignments, types.Assignment{
				ID:                assignID,
				Title:             assignData["attributes"].(map[string]interface{})["title"].(string),
				HTML:              assignData["attributes"].(map[string]interface{})["html"].(string),
				Done:              assignData["attributes"].(map[string]interface{})["done"].(bool),
				DueDateTime:       parseTime(assignData["attributes"].(map[string]interface{})["dueDateTime"].(string)),
				DeliverWorkOnline: assignData["attributes"].(map[string]interface{})["deliverWorkOnline"].(bool),
				OnlineDeliverURL:  assignData["attributes"].(map[string]interface{})["onlineDeliverUrl"].(string),
				Subject: types.Subject{
					ID:    subjectID,
					Label: subjectData["attributes"].(map[string]interface{})["label"].(string),
					Color: subjectData["attributes"].(map[string]interface{})["color"].(string),
				},
			})
		}

		result = append(result, types.TimetableDay{
			Date:        parseTime(day["attributes"].(map[string]interface{})["date"].(string)),
			Lessons:     lessons,
			Assignments: assignments,
		})
	}

	return result, nil
}

func parseTime(str string) time.Time {
	t, _ := time.Parse(time.RFC3339, str)
	return t
}
