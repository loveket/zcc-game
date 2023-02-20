package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const peoplepath = "./config/people.json"

func init() {
	NewPeoplePool()
}

var peopleSlice *PeopleSlice

type PeopleSlice struct {
	PeopleArr []*peopleMessage `json:"people_arr"`
}
type peopleMessage struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewPeoplePool() {
	file, err := os.Open(peoplepath)
	if err != nil {
		fmt.Println("file read err", err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("io read err", err)
	}
	ps := &PeopleSlice{}
	err = json.Unmarshal(data, &ps)
	//fmt.Println(ps, "***********")
	if err != nil {
		fmt.Println("Unmarshal err", err)
	}
	peopleSlice = ps
}

func GetPeopleName(id int) string {

	for _, v := range peopleSlice.PeopleArr {
		if v.Id == id {
			name := v.Name
			return name
		}
	}
	return ""
}
