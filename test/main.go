package main

import (
	"log"
	"os"
	"time"

	skolengo "github.com/espcaa/skolen-go"
)

func main() {

	// load the json data from "tokens.json"

	data, err := os.ReadFile("tokens.json")
	if err != nil {
		panic(err)
	}

	client, err := skolengo.NewClientFromJSON(data)

	if err != nil {
		panic(err)
	}
	// test the agenda

	log.Print("Fetching classes for user ", client.UserInfo.FullName, " at school ", client.School.Name)

	var yesterday_time = time.Now().AddDate(0, 0, -1)

	classes, err := client.GetTimetable(client.UserInfo.UserID, client.School.ID, client.School.EmsCode, yesterday_time, time.Now(), 100)

	if err != nil {
		panic(err)
	}

	for _, day := range classes {
		for _, class := range day.Lessons {
			println("  Class:", class.Subject.Label, "from", class.StartDateTime.Format("15:04"), "to", class.EndDateTime.Format("15:04"), "in room", class.Location)
		}
	}
}
