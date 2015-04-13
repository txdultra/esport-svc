package main

/*
import (
	"./libs"
	//"./models"
	"fmt"
	//"strconv"
	"testing"
)


func TestGetPersonal(t *testing.T) {
	l := &libs.Lives{}
	per := l.GetPersonal(1)
	if per == nil {
		t.Fatal("nil")
	}
	per = l.GetPersonal(1)
	fmt.Println(per)
}

func TestPersonals(t *testing.T) {
	l := &libs.Lives{}
	pers, err := l.GetPersonalLives(0, models.LIVE_STATUS_NOTHING, 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	pers, err = l.GetPersonalLives(2, models.LIVE_STATUS_LIVING, 1, 20)
	for _, v := range pers.List.([]models.LivePerson) {
		fmt.Println(v)
	}
}


func TestCreatePersonalLive(t *testing.T) {
	l := &libs.Lives{}
	pers, err := l.GetPersonalLives(1, 20, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("-------------before created-------------------")
	for _, v := range pers.List.([]models.LivePerson) {
		fmt.Println(v)
	}
	fmt.Println("-----------------------------------------------")

	per := &models.LivePerson{
		Name:       "老党在线",
		Uid:        2,
		ReptileUrl: "http://www.17173.com",
		Rep:        models.PER_REP_17173,
	}
	_, err = l.CreatePersonalLive(per, []int{2}, false)
	if err != nil {
		t.Fatal(err)
	}
	pers, err = l.GetPersonalLives(1, 20, nil)
	fmt.Println("-------------after created---------------------")
	for _, v := range pers.List.([]models.LivePerson) {
		fmt.Println(v)
	}
	fmt.Println("-----------------------------------------------")

	pers, err = l.GetPersonalLives(1, 20, nil)
	fmt.Println("-------------before updated gameid = 1---------------------")
	for _, v := range pers.List.([]models.LivePerson) {
		fmt.Println(v)
	}
	fmt.Println("-----------------------------------------------")

	l.UpdatePersonalLive(per, []int{1})
	pers, err = l.GetPersonalLives(1, 20, nil)
	fmt.Println("-------------after updated gameid = 1---------------------")
	for _, v := range pers.List.([]models.LivePerson) {
		fmt.Println(v)
	}
	fmt.Println("-----------------------------------------------")
}


func TestGetStreams(t *testing.T) {
	l := &libs.LiveOrgs{}
	lst := l.GetStreams(1)
	if len(lst) == 0 {
		fmt.Println("empty")
	}
	lst = l.GetStreams(1)
	for _, v := range lst {
		fmt.Println(v)
	}

	lts := l.GetChannels()
	for _, v := range lts {
		fmt.Println(v)
	}
}
*/
