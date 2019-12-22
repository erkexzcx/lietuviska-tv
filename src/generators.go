package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func updateDynamicChannels() {
	var wg sync.WaitGroup
	wg.Add(4) // 4 functions

	generateLRT(&wg)
	generateLRTPlius(&wg)
	generateLietuvosRytas(&wg)
	generateLnkGroup(&wg)

	wg.Wait()
}

func generateLRT(wg *sync.WaitGroup) {
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

	updateTVChannelURL("LRT HD", url)
	wg.Done()
}

func generateLRTPlius(wg *sync.WaitGroup) {
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

	updateTVChannelURL("LRT Plius HD", url)
	wg.Done()
}

func generateLietuvosRytas(wg *sync.WaitGroup) {
	lietuvosRytasURL, err := downloadContent("https://lib.lrytas.lt/geoip/get_token_live.php")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	updateTVChannelURL("Lietuvos rytas HD", string(lietuvosRytasURL))
	wg.Done()
}

func generateLnkGroup(wg *sync.WaitGroup) {
	// First, we need to download JSON from lnk api to see what is currently live:
	videosJSON, err := downloadContent("https://lnk.lt/api/main/live-page")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var result map[string]interface{}
	json.Unmarshal(videosJSON, &result)

	// List of videos. Usually TV show, or live TV
	OnlineVideos := result["videoGridCurrentLive"].(map[string]interface{})["videos"].([]interface{})
	OfflineVideos := result["videoGridNotLive"].(map[string]interface{})["videos"].([]interface{})

	parseVideos := func(videos *[]interface{}) {
		for _, v := range *videos {
			id := fmt.Sprintf("%v", v.(map[string]interface{})["id"])
			processLnkChannel(id)
		}
	}

	parseVideos(&OnlineVideos)
	parseVideos(&OfflineVideos)
	wg.Done()
}

func processLnkChannel(id string) {
	videoJSON, err := downloadContent("https://lnk.lt/api/main/video-page/xD/" + id + "/false")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(videoJSON, &result)

	params := result["videoConfig"].(map[string]interface{})["videoInfo"].(map[string]interface{})
	videoURL := params["videoUrl"].(string)
	videoURLTokens := params["secureTokenParams"].(string)

	switch videoURL {
	case "https://live.lnk.lt/lnk_live/tiesiogiai/playlist.m3u8":
		validateAndAddChannel("INFO TV HD", videoURL+videoURLTokens)
	case "https://live.lnk.lt/lnk_live/btv/playlist.m3u8":
		validateAndAddChannel("BTV HD", videoURL+videoURLTokens)
	case "https://live.lnk.lt/lnk_live/lnk/playlist.m3u8":
		validateAndAddChannel("LNK HD", videoURL+videoURLTokens)
	}
}

// IF TV channel is not valid - ignore it and do not update existing URLs
func validateAndAddChannel(title, link string) {
	res, err := http.Get(link)

	// If failed to perform HTTP request
	if err != nil {
		return
	}

	// If server did not return HTTP code 200
	if res.StatusCode != 200 {
		return
	}

	// If it doesn't contain something like "#EXT..."
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	if !strings.Contains(string(content), "\n#EXT") {
		return
	}

	// All good - update:
	updateTVChannelURL(title, link)
}

// downloadJSON downloads data. It's basically shortcut for GET request
func downloadContent(link string) ([]byte, error) {
	res, err := http.Get(link)
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
