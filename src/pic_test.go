package main

import (
	"fmt"
	"image"
	"log"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {

	reader, err := os.Open("5.6.mongodb.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	fmt.Println(image.Decode(reader))
	//fmt.Println(utils.ImageFormat(fd))

	return
}
