package main

import (
	"fmt"
	"github.com/c4urself/hackathon/mosaic"
	// "github.com/c4urself/hackathon/web"
)

func main() {
	fmt.Println(mosaic.MakeInstagramMosaic("pewdiepie", "/tmp/pewdiepie/photos", "/tmp/pewdiepie/audience", "/tmp/pewdiepie/mosaic/"))
	// mosaic.MakeMosaic("image.png", "./thumbnails/", "result.png", 10)
	// web.StartApp()
}
