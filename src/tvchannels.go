package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"sync"
)

type tvChannelsList map[string]*tvchannel

// Define all tv channels here with their static URLs (if they have):
var tvChannels = tvChannelsList{
	"TV3": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_303_1418375193.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV3_LT_HD/HLS_encr/TV3_LT_HD.m3u8",
		URLRoot: "https://cdn7.tvplayhome.lt/live/eds/TV3_LT_HD/HLS_encr/",
	},
	"LNK": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_301_1520339152.png",
		URL:     "",
		URLRoot: "",
	},
	"INFO TV": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_326_1467119944.png",
		URL:     "",
		URLRoot: "",
	},
	"LRT": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_306_1488445569.png",
		URL:     "",
		URLRoot: "",
	},
	"LRT Plius": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_307_1538382450.png",
		URL:     "",
		URLRoot: "",
	},
	"Lietuvos rytas": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_318_1539885851.png",
		URL:     "",
		URLRoot: "",
	},
}

var tvChannelsMutex = sync.Mutex{}

type tvchannel struct {
	Picture string
	URL     string
	URLRoot string
}

func (tvList tvChannelsList) renderPlaylist(w *http.ResponseWriter, addressHost string) {
	fmt.Fprintln(*w, "#EXTM3U")
	tvChannelsMutex.Lock()
	titles := make([]string, 0, len(tvList))
	for tvch := range tvList {
		titles = append(titles, tvch)
	}
	sort.Strings(titles)
	for _, title := range titles {
		fmt.Fprintf(*w, "#EXTINF:-1 tvg-logo=\"%s\", %s\n%s\n\n", tvList[title].Picture, title, "http://"+addressHost+"/channel/"+url.QueryEscape(title)+".m3u8")
	}
	tvChannelsMutex.Unlock()
}

var urlRootRe = regexp.MustCompile(`^(.+/)[^/]+$`)

func (tvList tvChannelsList) updateURL(title, url string) {
	match := urlRootRe.FindStringSubmatch(url)
	noEnding := url
	if match != nil {
		noEnding = match[1]
	}
	tvChannelsMutex.Lock()
	tvChannels[title].URL = url
	tvChannels[title].URLRoot = noEnding
	tvChannelsMutex.Unlock()
}
