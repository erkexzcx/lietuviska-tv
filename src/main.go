package main

import (
	"log"
	"net/http"
	"time"
)

func main() {

	log.Println("Starting...")

	// Initiate URLRoot (for static channels, before starting this app)
	initiateURLRoots()

	// Constantly update dynamic channels in the background
	go func() {
		for {
			go generateLietuvosRytas()
			go generateLRT()
			go generateLRTPlius()
			go generateLnkGroup()
			time.Sleep(1 * time.Hour)
		}
	}()

	http.HandleFunc("/iptv", renderPlaylist)
	http.HandleFunc("/iptv/", handleChannelRequest)

	log.Println("Started!")

	log.Fatal(http.ListenAndServe(":8989", nil))

}
