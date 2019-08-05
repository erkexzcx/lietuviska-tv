package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type tvlink struct {
	Title   string `json:"title"`
	Picture string `json:"picture"`
	URL     string `json:"url"`
}

var tvlinks = []tvlink{}

const fileFullPath = "links.json"

func main() {

	showTheAmazingError := func() {
		fmt.Fprintf(os.Stderr, "error: No argument was provided. Either use \"%v generate\" or \"%v show\"\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	args := os.Args
	if len(args) != 2 {
		showTheAmazingError()
	}
	switch os.Args[1] {
	case "generate":
		generate()
	case "show":
		show()
	default:
		showTheAmazingError()
	}

}

// generate regenerates and updates file
func generate() {

	loadFromFile()

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
		if title == "Žinios" {
			id := fmt.Sprintf("%v", el["id"])
			addEntry("LNK HD", "https://www.telia.lt/documents/20184/3686852/LNK-LOGO-HD.png", id)
		} else if title == "INFO TV HD kanalas internetu!" {
			id := fmt.Sprintf("%v", el["id"])
			addEntry("INFO TV HD", "https://www.telia.lt/documents/20184/3686852/INFO-LOGO-HD.png", id)
		}

	}

	// Lietuvos rytas TV hack as well:
	lietuvosRytasURL, err := downloadContent("https://lib.lrytas.lt/geoip/get_token_live.php")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	addEntry("Lietuvos rytas HD", "https://www.telia.lt/documents/20184/3686852/LRYTAS+TV+LOGOTIPAS.png", string(lietuvosRytasURL))

	saveToFile()

}

// show shows compiled m3u playlist from what is in the file
func show() {
	loadFromFile()

	for _, tv := range tvlinks {
		fmt.Println("#EXTM3U")
		fmt.Println("#EXTINF:-1 group-title=\"LT\" tvg-id=\"\" tvg-logo=\"" + tv.Picture + "\", " + tv.Title)
		fmt.Println(tv.URL)
		fmt.Println()
	}
}

// loadFromFile reads from file and save output to variable 'tvlinks'
func loadFromFile() {
	// Read file and attempt to parse previously known link
	content, err := ioutil.ReadFile(fileFullPath)
	if err == nil {
		err = json.Unmarshal(content, &tvlinks)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Define list of static live TVs
		tv3 := tvlink{
			Title:   "TV3 HD",
			Picture: "https://www.telia.lt/documents/20184/3686852/tv3-on-white.png",
			URL:     "https://cdn7.tvplayhome.lt/live/eds/TV3_LT_HD/HLS_encr/TV3_LT_HD.m3u8",
		}
		tvlinks = append(tvlinks, tv3)
	}
}

// saveToFile writes variable 'tvlinks' to file
func saveToFile() {
	// Write changes back to file:
	jsonTvLinks, _ := json.Marshal(tvlinks)
	f, err := os.Create(fileFullPath)
	check(err)
	defer f.Close()
	_, err = f.Write(jsonTvLinks)
	check(err)
	f.Sync()
}

// addEntry appends new entry to 'tvlinks' if it already exists (in terms of 'title' attribute)
func addEntry(title, picture, url string) {

	existsInArray := false
	for i, tvl := range tvlinks {
		if tvl.Title != title {
			continue
		}
		existsInArray = true
		// Update existing entry:
		tvlinks[i].Picture = picture
		tvlinks[i].URL = url
		break
	}

	if !existsInArray {
		// Add new entry
		tv := tvlink{
			Title:   title,
			Picture: picture,
			URL:     url,
		}
		tvlinks = append(tvlinks, tv)
	}
}

func processLnkChannel(title, picture, id string) {
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
	addEntry(title, picture, myURL)
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
