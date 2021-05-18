package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {

	http.HandleFunc("/watch", watch)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func watch(w http.ResponseWriter, r *http.Request) {
	yurl := fmt.Sprintf("https://youtube.com/watch?v=" + r.URL.Query().Get("v"))

	err := downloader(w, yurl)
	if err != nil {
		fmt.Fprintf(w, "cannot stream %v", err)
		return
	}
}

func downloader(w io.Writer, yurl string) error {

	readout, writein := io.Pipe()
	defer readout.Close()

	youtuber := exec.Command("youtube-dl", yurl, "-o-")

	youtuber.Stdout = writein
	youtuber.Stderr = os.Stderr
	ffmpeger := exec.Command("ffmpeg", "-i", "/dev/stdin", "-f", "mp3", "-ab", "96000", "-vn", "-")

	ffmpeger.Stdin = readout
	ffmpeger.Stdout = w
	ffmpeger.Stderr = os.Stderr

	go func() {
		if err := youtuber.Run(); err != nil {
			log.Fatalf("youtuber pipeline ist kaputt %v", err)
		}
	}()

	err := ffmpeger.Run()

	return err

}
