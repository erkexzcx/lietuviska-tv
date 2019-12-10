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

	// Update LNK group channels in the background
	go func() {
		for {
			go generateLietuvosRytas()
			go generateLRT()
			go generateLRTPlius()
			go generateLnkGroup()
			time.Sleep(15 * time.Minute)
		}
	}()

	// This provides a channels list, where URLs leads to the same server where you are hosting THIS application
	http.HandleFunc("/iptv", func(w http.ResponseWriter, r *http.Request) {
		tvChannels.renderPlaylist(&w, r.Host)
	})

	http.HandleFunc("/channel/", handleChannelRequest)
	http.HandleFunc("/link/", handleLinkRequest)

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

	tvChannels.updateURL("LRT", url)
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

	tvChannels.updateURL("LRT Plius", url)
}

func generateLietuvosRytas() {
	lietuvosRytasURL, err := downloadContent("https://lib.lrytas.lt/geoip/get_token_live.php")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	tvChannels.updateURL("Lietuvos rytas", string(lietuvosRytasURL))
}

func generateLnkGroup() {
	// First, we need to download JSON from lnk api to see what is currently live:
	videosJSON, err := downloadContent("https://lnk.lt/api/main/live-page")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Find IDs of videos :
	var result map[string]interface{}
	json.Unmarshal(videosJSON, &result)
	level1 := result["videoGridCurrentLive"].(map[string]interface{})
	level2 := level1["videos"].([]interface{})
	for _, v := range level2 {
		el := v.(map[string]interface{})

		title := fmt.Sprintf("%v", el["title"])
		if title == "Žinios" || title == "Labas vakaras, Lietuva" {
			id := fmt.Sprintf("%v", el["id"])
			processLnkChannel("LNK", id)
		} else if title == "INFO TV HD kanalas internetu!" {
			id := fmt.Sprintf("%v", el["id"])
			processLnkChannel("INFO TV", id)
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
	tvChannels.updateURL(title, url)

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
	return []byte(content), nil
}
