package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
)

type tvChannelsList map[string]*tvchannel

// Define all tv channels here with their static URLs (if they have):
var tvChannels = tvChannelsList{
	"TV3 HD": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_303_1575262629.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV3_LT_HD/HLS_encr/TV3_LT_HD.m3u8",
	},
	"TV3": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_303_1575262629.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV3_LT_SD/HLS_encr/TV3_LT_SD.m3u8",
	},
	"TV6": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_304_1572518420.jpg",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV6_LT_HD/HLS_encr/TV6_LT_HD.m3u8",
	},
	"TV8": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_545_1418375323.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV8_LT_SD/HLS_encr/TV8_LT_SD.m3u8",
	},
	"LRT": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_306_1488445569.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/LRT/HLS_encr/LRT.m3u8",
	},
	"LNK": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_301_1520339152.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/LNK/HLS_encr/LNK.m3u8",
	},
	"Lietuvos rytas": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_318_1539885851.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/Lietuvos_Ryto_TV/HLS_encr/Lietuvos_Ryto_TV.m3u8",
	},
	"BTV": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_305_1520956000.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/BTV/HLS_encr/BTV.m3u8",
	},
	"TV1": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_302_1421217703.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV1/HLS_encr/TV1.m3u8",
	},
	// "INFO TV": &tvchannel{ /* Actually it's LNK :D */
	// 	Picture: "https://cdn.tvstart.com/img/channel/logo_64_326_1467119944.png",
	// 	URL:     "https://cdn7.tvplayhome.lt/live/eds/InfoTV/HLS_encr/InfoTV.m3u8",
	// },
	"LRT Plius": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_307_1538382450.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/LRTKultura/HLS_encr/LRTKultura.m3u8",
	},
	"TV3 Film": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_303_1575262629.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV1000PremiumHD/HLS_encr/TV1000PremiumHD.m3u8",
	},
	"TV1000 Action": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_976_1490875823.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV1000_action_SD/HLS_encr/TV1000_action_SD.m3u8",
	},
	"TV1000 Ruskoje kino": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_978_1490875981.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV1000_rus_kino_SD/HLS_encr/TV1000_rus_kino_SD.m3u8",
	},
	"FOX LT": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_908_1454493659.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/FOX_LT/HLS_encr/FOX_LT.m3u8",
	},
	"Sony": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_957_1490876156.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/Sony_Channel/HLS_encr/Sony_Channel.m3u8",
	},
	"TV1000": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_977_1490875718.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TV1000_SD/HLS_encr/TV1000_SD.m3u8",
	},
	"Sony Turbo": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_70_1470387888.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/Sony_Turbo/HLS_encr/Sony_Turbo.m3u8",
	},
	"Viasat Explore": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_972_1420814929.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/ViasatExploreHD/HLS_encr/ViasatExploreHD.m3u8",
	},
	"Viasat History": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_973_1420814953.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/ViasatHistoryHD/HLS_encr/ViasatHistoryHD.m3u8",
	},
	"Viasat Nature": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_958_1420814982.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/Viasat_Nature_SD/HLS_encr/Viasat_Nature_SD.m3u8",
	},
	"NatioNAL Geographic": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_936_1454495068.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/NatGeoHD/HLS_encr/NatGeoHD.m3u8",
	},
	"Ohota y Ribalka": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_145_1454423846.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/Ohota_i_Ribalka/HLS_encr/Ohota_i_Ribalka.m3u8",
	},
	"BBC World News": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_80_1472213898.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/BBC/HLS_encr/BBC.m3u8",
	},
	"Euronews": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_79_1491386244.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/Euronews_RUS/HLS_encr/Euronews_RUS.m3u8",
	},
	"MTV Hits": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_69_1491913231.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/MTV_Hits/HLS_encr/MTV_Hits.m3u8",
	},
	"My Hits": &tvchannel{
		Picture: "https://r-scale-88.static.go3.tv/scale/go3/webuploads/rest/upload/live/1179183/images/6376570?dsth=512&dstw=512&srcmode=0&srcx=0&srcy=0&quality=65&type=0&srcw=1/1&srch=1/1",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/My_Hits/HLS_encr/My_Hits.m3u8",
	},
	"PBK LT": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_312_1397718439.jpg",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/PBK_LT/HLS_encr/PBK_LT.m3u8",
	},
	"Ren TV Baltic LT": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_199_1490875230.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/Ren_TV_LT/HLS_encr/Ren_TV_LT.m3u8",
	},
	"NTV Mir Baltic LT": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_106_1454492760.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/NTV_Mir_Baltic_LT/HLS_encr/NTV_Mir_Baltic_LT.m3u8",
	},
	"CTC Baltija": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_546_1490875440.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/CTCBaltics/HLS_encr/CTCBaltics.m3u8",
	},
	"TV3 Sport": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_919_1575344906.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/ViasatSportBaltics_HD/HLS_encr/ViasatSportBaltics_HD.m3u8",
	},
	"TV3 Sport 2": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_1101_1575344950.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/TVPlaySportPlusHD/HLS_encr/TVPlaySportPlusHD.m3u8",
	},
	"Fight sport": &tvchannel{
		Picture: "https://r-scale-d9.dcs.redlabs.pl/scale/AMB/webuploads/rest-uat/upload/live/939682/images/4502626?dsth=512&dstw=512&srcmode=0&srcx=0&srcy=0&quality=65&type=0&srcw=1/1&srch=1/1",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/FightSportHD/HLS_encr/FightSportHD.m3u8",
	},
	"Setanta Sports": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_933_1454422589.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/Setanta/HLS_encr/Setanta.m3u8",
	},
	"NBA": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_338_1454420591.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/NBA_TV_HD/HLS_encr/NBA_TV_HD.m3u8",
	},
	"Sport 1": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_336_1467180549.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/Sport1/HLS_encr/Sport1.m3u8",
	},
	"Esports": &tvchannel{
		Picture: "https://r-scale-b4.static.go3.tv/scale/go3-test/webuploads/rest-uat/upload/live/1048227/images/5538636?dsth=512&dstw=512&srcmode=0&srcx=0&srcy=0&quality=65&type=0&srcw=1/1&srch=1/1",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/EsportHD/HLS_encr/EsportHD.m3u8",
	},
	"Nick Junior": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_560_1454415832.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/Nick_JR/HLS_encr/Nick_JR.m3u8",
	},
	"Kidzone Plus": &tvchannel{
		Picture: "https://r-scale-52.static.go3.tv/scale/go3/webuploads/rest/upload/live/1218529/images/6615837?dsth=512&dstw=512&srcmode=0&quality=65&type=0&srcx=1&srcy=1&srcw=1/1&srch=1/1",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/KidZone_Plus/HLS_encr/KidZone_Plus.m3u8",
	},
	"Nickelodeon LT": &tvchannel{
		Picture: "https://cdn.tvstart.com/img/channel/logo_64_960_1418121302.png",
		URL:     "https://cdn7.tvplayhome.lt/live/eds/Nickelodeon/HLS_encr/Nickelodeon.m3u8",
	},
	/* Dynamic channels */
	"LNK HD":            &tvchannel{Picture: "https://cdn.tvstart.com/img/channel/logo_64_301_1520339152.png"},
	"INFO TV HD":        &tvchannel{Picture: "https://cdn.tvstart.com/img/channel/logo_64_326_1467119944.png"},
	"BTV HD":            &tvchannel{Picture: "https://cdn.tvstart.com/img/channel/logo_64_305_1520956000.png"},
	"LRT HD":            &tvchannel{Picture: "https://cdn.tvstart.com/img/channel/logo_64_306_1488445569.png"},
	"LRT Plius HD":      &tvchannel{Picture: "https://cdn.tvstart.com/img/channel/logo_64_307_1538382450.png"},
	"Lietuvos rytas HD": &tvchannel{Picture: "https://cdn.tvstart.com/img/channel/logo_64_318_1539885851.png"},
}

var tvMutex = sync.RWMutex{}

type tvchannel struct {
	Picture string
	URL     string
	URLRoot string
}

func renderPlaylist(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "#EXTM3U")

	tvMutex.RLock()
	titles := make([]string, 0, len(tvChannels))
	for tvch := range tvChannels {
		titles = append(titles, tvch)
	}
	sort.Strings(titles)
	for _, title := range titles {
		fmt.Fprintf(w, "#EXTINF:-1 tvg-logo=\"%s\", %s\n%s\n\n", tvChannels[title].Picture, title, "http://"+r.Host+"/iptv/"+url.QueryEscape(title)+".m3u8")
	}
	tvMutex.RUnlock()
}

func initiateURLRoots() {
	// Mutex is not needed, since we do it on startup
	for _, tvc := range tvChannels {
		if tvc.URL == "" {
			continue
		}
		tvc.URLRoot = deleteAfterLastSlash(&(tvc.URL))
	}
}

func updateTVChannelURL(title, url string) {
	URLRootNew := deleteAfterLastSlash(&url)
	tvMutex.Lock()
	tvChannels[title].URL = url
	tvChannels[title].URLRoot = URLRootNew
	tvMutex.Unlock()
}

func deleteAfterLastSlash(str *string) string {
	return (*str)[0 : strings.LastIndex(*str, "/")+1]
}
