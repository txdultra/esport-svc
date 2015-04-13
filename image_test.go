package main

import (
	"github.com/disintegration/imaging"
	"io/ioutil"
	"libs"
	"testing"
)

func TestImageResize(t *testing.T) {
	file, err := ioutil.ReadFile("f:\\146.jpg")
	if err != nil {
		t.Fatal(err)
	}
	//nga, _ := img.Resize(file, 0, 100)
	//nga, _ := img.Fit(file, 100, 100
	img, _ := libs.NewByteImage(file)
	nga := img.CropCenter(200, 200)
	imaging.Save(nga, "f:\\eeee.jpg")
}
