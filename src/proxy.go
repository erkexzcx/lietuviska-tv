package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func print404(w http.ResponseWriter, customMessage interface{}) {
	log.Println(customMessage)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 page not found"))
}

var urlChannelRe = regexp.MustCompile(`^\/channel\/(.+).m3u8$`)

func handleFirstm3u8(w http.ResponseWriter, r *http.Request) {

	match := urlChannelRe.FindStringSubmatch(r.URL.Path)
	if match == nil {
		print404(w, "Unable to properly extract data from request '"+r.URL.Path+"'!")
		return
	}

	encodedChannelName := match[1]

	decodedChannelName, err := url.QueryUnescape(encodedChannelName)
	if err != nil {
		print404(w, "Unable to decode channel '"+encodedChannelName+"'!")
		return
	}

	tvChannelsMutex.Lock()
	el, ok := tvChannels[decodedChannelName]
	tvChannelsMutex.Unlock()
	if !ok {
		print404(w, "Unable to find channel '"+decodedChannelName+"'!")
		return
	}

	tvChannelsMutex.Lock()
	requiredURL := el.URL
	tvChannelsMutex.Unlock()

	if requiredURL == "" {
		print404(w, "Channel '"+decodedChannelName+"' does not have URL assigned!")
		return
	}

	resp, err := http.Get(requiredURL)
	if err != nil {
		print404(w, err)
		return
	}

	// Print line by line, and if it doesn't start with '#' - append URL for next step:
	w.WriteHeader(http.StatusOK)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			line = "http://" + r.Host + "/channel2/" + encodedChannelName + "/" + line
		}
		w.Write([]byte(line + "\n"))
	}
}

var urlChannel2Re = regexp.MustCompile(`^\/channel2\/([^\/]+)\/(.+)$`)

func handleSecondm3u8(w http.ResponseWriter, r *http.Request) {

	match := urlChannel2Re.FindStringSubmatch(r.URL.Path)

	if match == nil {
		print404(w, "Unable to properly extract data from request '"+r.URL.Path+"'!")
		return
	}

	encodedChannelName := match[1]
	relativePath := match[2]

	decodedChannelName, err := url.QueryUnescape(encodedChannelName)
	if err != nil {
		print404(w, "Unable to decode channel '"+encodedChannelName+"'!")
		return
	}

	tvChannelsMutex.Lock()
	el, ok := tvChannels[decodedChannelName]
	tvChannelsMutex.Unlock()
	if !ok {
		print404(w, "Unable to find channel '"+decodedChannelName+"'!")
		return
	}

	tvChannelsMutex.Lock()
	requiredURL := el.URLRoot
	tvChannelsMutex.Unlock()

	if requiredURL == "" {
		print404(w, "Channel '"+decodedChannelName+"' does not have root URL assigned!")
		return
	}

	resp, err := http.Get(requiredURL + relativePath)
	if err != nil {
		print404(w, err)
		return
	}

	// Print line by line, and if it doesn't start with '#' - append URL for next step:
	w.WriteHeader(http.StatusOK)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			line = "http://" + r.Host + "/channel3/" + encodedChannelName + "/" + line
		}
		w.Write([]byte(line + "\n"))
	}
}

var urlChannel3Re = regexp.MustCompile(`^\/channel3\/([^\/]+)\/(.+)$`)

func handleThirdm3u8(w http.ResponseWriter, r *http.Request) {

	match := urlChannel3Re.FindStringSubmatch(r.URL.Path)
	if match == nil {
		print404(w, "Unable to properly extract data from request '"+r.URL.Path+"'!")
		return
	}

	encodedChannelName := match[1]
	relativePath := match[2]

	decodedChannelName, err := url.QueryUnescape(encodedChannelName)
	if err != nil {
		print404(w, "Unable to decode channel '"+encodedChannelName+"'!")
		return
	}

	tvChannelsMutex.Lock()
	el, ok := tvChannels[decodedChannelName]
	tvChannelsMutex.Unlock()
	if !ok {
		print404(w, "Unable to find channel '"+decodedChannelName+"'!")
		return
	}

	tvChannelsMutex.Lock()
	requiredURL := el.URLRoot
	tvChannelsMutex.Unlock()

	if requiredURL == "" {
		print404(w, "Channel '"+decodedChannelName+"' does not have root URL assigned!")
		return
	}

	resp, err := http.Get(requiredURL + relativePath)
	if err != nil {
		print404(w, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print404(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
