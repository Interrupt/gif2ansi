package main

import (
	"fmt"
	"image/gif"
	"log"
	"os"
	"path"
	"time"

	"atomicgo.dev/cursor"
)

var imageFile string
var border int

func main() {
	if len(os.Args) != 2 {
		usage := `Usage:
	%s <path to gif>`
		fmt.Println(fmt.Sprintf(usage, path.Base(os.Args[0])))
		os.Exit(1)
	}

	// Top gutter
	fmt.Println("")

	imageFile := os.Args[1]
	f, err := os.Open(imageFile)
	if err != nil {
		log.Fatalf("Error opening gif image [%s]: %v", imageFile, err)
	}
	defer f.Close()

	gif, err := gif.DecodeAll(f)
	if err != nil {
		log.Fatalf("Error parsing gif image: %v", err)
	}

	frameCount := len(gif.Image)

	// Draw each gif frame into the terminal
	for i := 0; i < frameCount; i++ {
		renderStart := time.Now()

		ansi, err := RenderTrueColor(gif.Image[i])
		if err != nil {
			log.Fatalf("Error converting gif frame to ANSI: %v", err)
		}
		fmt.Print(ansi)

		if i+1 != frameCount {
			cursor.Up((gif.Config.Height + 1) / 2)
		}

		// Remove render time from the sleep, to speed up drawing
		renderDuration := int(time.Since(renderStart).Milliseconds())

		// Delay is in 100s of a second
		time.Sleep(time.Duration((gif.Delay[i]*10)-renderDuration) * time.Millisecond)
	}

	// Bottom gutter
	fmt.Println("")
}
