package main

/*
import (
	"./libs"
	"./models"
	"fmt"
	//"strconv"
	"testing"
)


func TestCreatePlaylist(t *testing.T) {
	vods := &libs.Vods{}
	pls := models.VideoPlaylist{}
	pls.Title = "草泥马"
	pls.Uid = 1
	pls.Img = 1
	id, err := vods.CreatePlaylist(pls, map[int64]int{
		9:  1,
		10: 2,
		11: 3,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("new playlist id:" + strconv.FormatInt(id, 10))
}

func TestGetPlaylist(t *testing.T) {
	vods := &libs.Vods{}
	pls := vods.GetPlaylist(1)
	if pls == nil {
		t.Fatal("错误")
	}
	pls2 := vods.GetPlaylist(1)
	fmt.Println(pls2)
}

func TestGetPlaylistVods(t *testing.T) {
	vods := &libs.Vods{}
	_, err := vods.GetPlaylistVods(1, 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	plvods2, _ := vods.GetPlaylistVods(1, 1, 20)
	fmt.Println(plvods2)
}

func TestUpdatePLsVodNos(t *testing.T) {
	vods := &libs.Vods{}
	err := vods.UpdatePLsVodNos(1, map[int64]int{
		9:  2,
		10: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	plvods2, _ := vods.GetPlaylistVods(1, 1, 20)
	for _, v := range plvods2.List.([]*models.Video) {
		fmt.Println(v)
	}
}


func TestRemovePlsVod(t *testing.T) {
	vods := &libs.Vods{}
	err := vods.RemovePlsVod(1, 9)
	if err != nil {
		t.Fatal(err)
	}
	plvods2, _ := vods.GetPlaylistVods(1, 1, 20)
	for _, v := range plvods2.List.([]*models.Video) {
		fmt.Println(v)
	}
}


func TestAppenedPlsVods(t *testing.T) {
	vods := &libs.Vods{}
	err := vods.AppenedPlsVods(1, map[int64]int{
		13: 2,
		14: 10,
		10: 13,
	})
	if err != nil {
		t.Fatal(err)
	}
	plvods2, _ := vods.GetPlaylistVods(1, 1, 2)
	for _, v := range plvods2.List.([]*models.Video) {
		fmt.Println(v)
	}
}

*/
