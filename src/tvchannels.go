package main

import (
	"fmt"
	"net/http"
	"sync"
)

type tvChannelsList map[string]*tvchannel

// Define all tv channels here. Code will keep updating URLs, while web
// server won't show these tv channels without URL.
var tvChannels = tvChannelsList{
	"TV3": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_303_1418375193.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV3_LT_HD/HLS_encr/TV3_LT_HD.m3u8",
	},
	"LNK": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_301_1520339152.png",
		URL:     "",
	},
	"INFO TV": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_326_1467119944.png",
		URL:     "",
	},
	"LRT": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_306_1488445569.png",
		URL:     "",
	},
	"LRT Plius": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_307_1538382450.png",
		URL:     "",
	},
	"Lietuvos rytas": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_318_1539885851.png",
		URL:     "",
	},
}

var tvChannelsMutex = sync.Mutex{}

type tvchannel struct {
	Picture string
	URL     string
}

func (tvList tvChannelsList) renderPlaylist(w *http.ResponseWriter) {

	// Some channels are always available, just direct link is needed to be parsed
	var wg sync.WaitGroup
	wg.Add(3)

	go generateLietuvosRytas(&wg)
	go generateLRT(&wg)
	go generateLRTPlius(&wg)

	wg.Wait()

	fmt.Fprintln(*w, "#EXTM3U")
	tvChannelsMutex.Lock()
	for title, channel := range tvList {
		if channel.URL == "" {
			continue
		}
		fmt.Fprintf(*w, "#EXTINF:-1 tvg-logo=\"%s\", %s\n%s\n\n", channel.Picture, title, channel.URL)
	}
	tvChannelsMutex.Unlock()
}

func (tvList tvChannelsList) updateURL(title, url string) {
	tvChannelsMutex.Lock()
	tvChannels[title].URL = url
	tvChannelsMutex.Unlock()
}
