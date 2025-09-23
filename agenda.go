package skolengo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/espcaa/skolen-go/types"
	utils "github.com/espcaa/skolen-go/utils"
)

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

type BaseResponse struct {
	Data     []map[string]interface{} `json:"data"`
	Included []map[string]interface{} `json:"included"`
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok && v != nil {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func parseTime(str string) time.Time {
	t, _ := time.Parse(time.RFC3339, str)
	return t
}

func (c *Client) GetTimetable(userID, schoolID, emsCode string, periodStart, periodEnd time.Time, limit int) ([]types.TimetableDay, error) {
	// Ensure we have a valid token before making the API call
	if err := c.EnsureValidToken(); err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

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

	// Map all included items for quick lookup
	includedMap := make(map[string]map[string]interface{})
	for _, item := range baseResp.Included {
		id := getString(item, "id")
		typ := getString(item, "type")
		if id != "" && typ != "" {
			includedMap[fmt.Sprintf("%s:%s", typ, id)] = item
		}
	}

	var result []types.TimetableDay
	for _, day := range baseResp.Data {
		dayAttr, ok := day["attributes"].(map[string]interface{})
		if !ok {
			continue
		}

		relationships, ok := day["relationships"].(map[string]interface{})
		if !ok {
			continue
		}

		rawLessons := utils.GetMultipleRelations(relationships["lessons"].(map[string]interface{}))
		rawAssignments := utils.GetMultipleRelations(relationships["homeworkAssignments"].(map[string]interface{}))

		var lessons []types.Lesson
		for _, rawLesson := range rawLessons {
			lessonID := getString(rawLesson, "id")
			lessonData, ok := includedMap["lesson:"+lessonID]
			if !ok {
				continue
			}

			attr, ok := lessonData["attributes"].(map[string]interface{})
			if !ok {
				continue
			}

			subjectRel := utils.GetSingleRelation(lessonData["relationships"].(map[string]interface{})["subject"].(map[string]interface{}))
			subject := types.Subject{}
			if subjectRel != nil {
				subjID := getString(subjectRel, "id")
				if subjData, ok := includedMap["subject:"+subjID]; ok {
					subjAttr := subjData["attributes"].(map[string]interface{})
					subject = types.Subject{
						ID:    subjID,
						Label: getString(subjAttr, "label"),
						Color: getString(subjAttr, "color"),
					}
				}
			}

			var teachers []types.Teacher
			for _, t := range utils.GetMultipleRelations(lessonData["relationships"].(map[string]interface{})["teachers"].(map[string]interface{})) {
				tID := getString(t, "id")
				if tData, ok := includedMap["teacher:"+tID]; ok {
					tAttr := tData["attributes"].(map[string]interface{})
					teachers = append(teachers, types.Teacher{
						ID:        tID,
						Title:     getString(tAttr, "title"),
						FirstName: getString(tAttr, "firstName"),
						LastName:  getString(tAttr, "lastName"),
						PhotoURL:  getString(tAttr, "photoUrl"),
					})
				}
			}

			lessons = append(lessons, types.Lesson{
				ID:                          lessonID,
				StartDateTime:               parseTime(getString(attr, "startDateTime")),
				EndDateTime:                 parseTime(getString(attr, "endDateTime")),
				Location:                    getString(attr, "location"),
				Canceled:                    getBool(attr, "canceled"),
				AnyContent:                  getBool(attr, "anyContent"),
				AnyHomeworkToDoForTheLesson: getBool(attr, "anyHomeworkToDoForTheLesson"),
				AnyHomeworkToDoAfterLesson:  getBool(attr, "anyHomeworkToDoAfterLesson"),
				Subject:                     subject,
				Teachers:                    teachers,
			})
		}

		var assignments []types.Assignment
		for _, rawAssign := range rawAssignments {
			assignID := getString(rawAssign, "id")
			assignData, ok := includedMap["homework:"+assignID]
			if !ok {
				continue
			}

			attr, ok := assignData["attributes"].(map[string]interface{})
			if !ok {
				continue
			}

			subjectRel := utils.GetSingleRelation(assignData["relationships"].(map[string]interface{})["subject"].(map[string]interface{}))
			subject := types.Subject{}
			if subjectRel != nil {
				subjID := getString(subjectRel, "id")
				if subjData, ok := includedMap["subject:"+subjID]; ok {
					subjAttr := subjData["attributes"].(map[string]interface{})
					subject = types.Subject{
						ID:    subjID,
						Label: getString(subjAttr, "label"),
						Color: getString(subjAttr, "color"),
					}
				}
			}

			assignments = append(assignments, types.Assignment{
				ID:                assignID,
				Title:             getString(attr, "title"),
				HTML:              getString(attr, "html"),
				Done:              getBool(attr, "done"),
				DueDateTime:       parseTime(getString(attr, "dueDateTime")),
				DeliverWorkOnline: getBool(attr, "deliverWorkOnline"),
				OnlineDeliverURL:  getString(attr, "onlineDeliverUrl"),
				Subject:           subject,
			})
		}

		result = append(result, types.TimetableDay{
			Date:        parseTime(getString(dayAttr, "date")),
			Lessons:     lessons,
			Assignments: assignments,
		})
	}

	return result, nil
}
