package main

import (
	"fmt"
	"image/gif"
	"log"
	"os"
	"path"
	"time"

	"atomicgo.dev/cursor"
	"github.com/mattn/go-tty"
)

var imageFile string
var border int
var height int

var ttyTerm *tty.TTY
var err error

func main() {
	if len(os.Args) != 2 {
		usage := `Usage:
	%s <path to gif>`
		fmt.Println(fmt.Sprintf(usage, path.Base(os.Args[0])))
		os.Exit(1)
	}

	ttyTerm, err = tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer ttyTerm.Close()

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
	height = gif.Config.Height + 1

	go watchForKeyPress()

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

func watchForKeyPress() {
	// Wait for a keypress
	_, err := ttyTerm.ReadRune()
	if err != nil {
		log.Fatal(err)
	}

	// Move the cursor down past the gif manually
	for i := 0; i < (height/2)+1; i++ {
		fmt.Printf("\n")
	}

	ttyTerm.Close()
	os.Exit(1)
}
