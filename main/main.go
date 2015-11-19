package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"

	"github.com/tonyhb/goicns"
)

var (
	flOut = flag.String("o", "", "Output filename")
)

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Usage: goicns -o IconSet.icns image.png")
		return
	}

	path := flag.Args()[0]

	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println("Error opening image %s: %s", path, err.Error())
		return
	}

	img, err := png.Decode(f)
	if err != nil {
		fmt.Println("Error decoding image %s: %s", path, err.Error())
		return
	}

	icns := goicns.NewICNS(img)
	if err = icns.Construct(); err != nil {
		fmt.Println("Error encofing ICNS %s: %s", path, err.Error())
		return
	}

	icns.WriteToFile(*flOut, 0666)
}
