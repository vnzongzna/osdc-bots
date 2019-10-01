package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

var finallist string //variable to store final list of meetups

//-------structs to parse json data from Meetup API
type venue struct {
	Venue string `json:"name"`
}
type group struct {
	GroupName string `json:"name"`
}
type meetuplist struct {
	Name  string `json:"name"`
	Venue venue  `json: "venue"`
	Date  string `json:"local_date"`
	Time  string `json:"local_time"`
	Group group  `json: "group"`
	Link  string `json:"link"`
}

//struct to keep track of OSDC meetups using a JSON file
type nextmeetup struct {
	Name string `json: "Name"`
	Date string `json: "Date"`
}

//fetching details of meetup of the group's url
func getMeetups(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var meetups []meetuplist
	json.Unmarshal([]byte(body), &meetups)
	if (len(meetups)) > 0 {
		finallist = finallist + "Title -" + "\t" + meetups[0].Name + "\n" + "Date -" + "\t" + meetups[0].Date + "\n" + "Link -" + "\t" + meetups[0].Link + "\n\n"
	}
}

func addmeetup(ID int64, msgtext string) {
	args := strings.Fields(msgtext)
	data := nextmeetup{
		Name: args[1],
		Date: args[2],
	}
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("meetups.json", file, 0644)
	bot.Send(tbot.NewMessage(ID, "Meetup added successfully."))
}
func getnextmeetup(ID int64) {
	file, _ := ioutil.ReadFile("meetups.json")
	data := nextmeetup{}
	_ = json.Unmarshal([]byte(file), &data)
	fmt.Println(data.Name)
	nxtmeetupdata := "Title -" + "\t" + data.Name + "\n" + "Date -" + "\t" + data.Date
	bot.Send(tbot.NewMessage(ID, nxtmeetupdata))

}
