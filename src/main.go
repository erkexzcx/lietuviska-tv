package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
			time.Sleep(2 * time.Hour)
		}
	}()

	http.HandleFunc("/iptv", renderPlaylist)
	http.HandleFunc("/iptv/", handleChannelRequest)

	log.Println("Started!")

	log.Fatal(http.ListenAndServe(":8989", nil))

}

func generateLRT() {
	ltvURL, err := downloadContent("https://www.lrt.lt/servisai/stream_url/live/get_live_url.php?channel=LTV1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(ltvURL), &result)
	level1 := result["response"].(map[string]interface{})
	level2 := level1["data"].(map[string]interface{})
	url := fmt.Sprintf("%v", level2["content"])

	updateTVChannelURL("LRT HD (D)", url)
}

func generateLRTPlius() {
	ltvURL, err := downloadContent("https://www.lrt.lt/servisai/stream_url/live/get_live_url.php?channel=LTV2")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(ltvURL), &result)
	level1 := result["response"].(map[string]interface{})
	level2 := level1["data"].(map[string]interface{})
	url := fmt.Sprintf("%v", level2["content"])

	updateTVChannelURL("LRT Plius (D)", url)
}

func generateLietuvosRytas() {
	lietuvosRytasURL, err := downloadContent("https://lib.lrytas.lt/geoip/get_token_live.php")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	updateTVChannelURL("Lietuvos rytas (D)", string(lietuvosRytasURL))
}

func generateLnkGroup() {
	// First, we need to download JSON from lnk api to see what is currently live:
	videosJSON, err := downloadContent("https://lnk.lt/api/main/live-page")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Find IDs of videos:
	var result map[string]interface{}
	json.Unmarshal(videosJSON, &result)
	level1 := result["videoGridNotLive"].(map[string]interface{})
	level2 := level1["videos"].([]interface{})
	for _, v := range level2 {
		el := v.(map[string]interface{})

		title := fmt.Sprintf("%v", el["title"])
		if title == "LNK HD kanalas internetu!" {
			id := fmt.Sprintf("%v", el["id"])
			processLnkChannel("LNK HD (D)", id)
		} else if title == "INFO TV HD kanalas internetu!" {
			id := fmt.Sprintf("%v", el["id"])
			processLnkChannel("INFO TV (D)", id)
		}

	}
}

func processLnkChannel(title, id string) {
	// download another JSON
	videoJSON, err := downloadContent("https://lnk.lt/api/main/video-page/xD/" + id + "/false")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var result map[string]interface{}
	json.Unmarshal(videoJSON, &result)
	level1 := result["videoConfig"].(map[string]interface{})
	level2 := level1["videoInfo"].(map[string]interface{})

	url := fmt.Sprintf("%v%v", level2["videoUrl"], level2["secureTokenParams"])
	updateTVChannelURL(title, url)

}

// downloadJSON downloads data. It's basically shortcut for GET request
func downloadContent(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
