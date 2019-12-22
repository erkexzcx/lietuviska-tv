package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func main() {

	log.Println("Starting...")

	// Initiate URLRoot (for static channels, before starting this app)
	initiateURLRoots()

	// Constantly update dynamic channels in the background
	updateDynamicChannels()
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			go updateDynamicChannels()
		}
	}()

	http.HandleFunc("/iptv", renderPlaylist)
	http.HandleFunc("/iptv/", handleChannelRequest)

	log.Println("Started!")

	ips, err := getAvailableURLs()
	if err == nil {
		fmt.Println("\nUse below IP addresses to reach M3U playlist:")
		for _, link := range ips {
			fmt.Println("\t" + link)
		}
		fmt.Println()
	}

	log.Fatal(http.ListenAndServe(":8989", nil))

}

func getAvailableURLs() ([]string, error) {
	var IPs []string
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, v := range addresses {
		address := v.String()
		if strings.Contains(address, "::") {
			continue
		}
		IPs = append(IPs, "http://"+strings.SplitN(address, "/", 2)[0]+":8989/iptv")
	}
	return IPs, nil
}
