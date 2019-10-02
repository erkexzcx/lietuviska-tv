package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Define all channels here. Code will keep updating
// URLs, while web server won't show these channels
// without URL.
var tvlist = map[string]tvchannel{
	"TV3": tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_303_1418375193.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV3_LT_HD/HLS_encr/TV3_LT_HD.m3u8",
	},
	"LNK": tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_301_1520339152.png",
		URL:     "",
	},
	"INFO TV": tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_326_1467119944.png",
		URL:     "",
	},
	"LRT": tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_306_1488445569.png",
		URL:     "",
	},
	"LRT Plius": tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_307_1538382450.png",
		URL:     "",
	},
	"Lietuvos rytas": tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_318_1539885851.png",
		URL:     "",
	},
}

var tvlistMutex = &sync.Mutex{}

type tvchannel struct {
	Picture string `json:"picture"`
	URL     string `json:"url"`
}

func main() {

	// Run linkUpdater in the background
	go linkUpdater()

	http.HandleFunc("/iptv", func(w http.ResponseWriter, r *http.Request) {
		renderPlaylist(&w)
	})

	http.HandleFunc("/epg", func(w http.ResponseWriter, r *http.Request) {
		// TODO
	})

	log.Fatal(http.ListenAndServe(":8989", nil))

}

func renderPlaylist(w *http.ResponseWriter) {

	// Some channels are always available, just direct link is needed to be parsed
	var wg sync.WaitGroup
	wg.Add(3)

	go generateLietuvosRytas(&wg)
	go generateLRT(&wg)
	go generateLRTPlius(&wg)

	wg.Wait()

	fmt.Fprintln(*w, "#EXTM3U")
	tvlistMutex.Lock()
	for name, channel := range tvlist {
		if channel.URL == "" {
			continue
		}
		fmt.Fprintf(*w, "#EXTINF:-1 tvg-logo=\"%s\", %s\n%s\n\n", channel.Picture, name, channel.URL)
	}
	tvlistMutex.Unlock()
}

func linkUpdater() {
	for {
		generateLnkGroup()
		time.Sleep(10 * time.Minute)
	}
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

	tvlistMutex.Lock()
	x := tvlist["LRT"]
	x.URL = url
	tvlist["LRT"] = x
	tvlistMutex.Unlock()

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

	tvlistMutex.Lock()
	x := tvlist["LRT Plius"]
	x.URL = url
	tvlist["LRT Plius"] = x
	tvlistMutex.Unlock()

	wg.Done()
}

func generateLietuvosRytas(wg *sync.WaitGroup) {
	lietuvosRytasURL, err := downloadContent("https://lib.lrytas.lt/geoip/get_token_live.php")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	tvlistMutex.Lock()
	x := tvlist["Lietuvos rytas"]
	x.URL = string(lietuvosRytasURL)
	tvlist["Lietuvos rytas"] = x
	tvlistMutex.Unlock()

	wg.Done()
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

	myURL := fmt.Sprintf("%v%v", level2["videoUrl"], level2["secureTokenParams"])

	tvlistMutex.Lock()
	x := tvlist[title]
	x.URL = myURL
	tvlist[title] = x
	tvlistMutex.Unlock()

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

// check is simpliefied one line check for file IO operations
func check(e error) {
	if e != nil {
		panic(e)
	}
}
