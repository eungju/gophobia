package main

import (
	"fmt"
	"github.com/rjeczalik/notify"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	runTest()

	events := make(chan notify.EventInfo, 1)
	if err := notify.Watch("./...", events, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(events)

	lastT := time.Unix(0, 0)
	for {
		select {
		case event := <-events:
			relPath, err := filepath.Rel(baseDir, event.Path())
			if err != nil {
				log.Fatal(err)
			}
			if !strings.HasPrefix(relPath, ".") {
				log.Println(event)
				d := time.Now().Sub(lastT).Seconds()
				if d >= 3.0 {
					runTest()
				}
				lastT = time.Now()
			}
		}
	}
}

func runTest() {
	cmd := exec.Command("go", "test", "./...")
	outerr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", outerr)
		log.Println(err)
	} else {
		log.Println("ok.")
	}
}
