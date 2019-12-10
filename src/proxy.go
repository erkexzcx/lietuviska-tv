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

var regexChannelWithURL = regexp.MustCompile(`^\/\w+\/([^\/]+)\/(.+)$`)

var regexChannelOnly = regexp.MustCompile(`^\/channel\/([^\/]+)\.m3u8$`)

func handleChannelRequest(w http.ResponseWriter, r *http.Request) {

	match := regexChannelOnly.FindStringSubmatch(r.URL.Path)
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

	w.WriteHeader(http.StatusOK)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			line = "http://" + r.Host + "/link/" + encodedChannelName + "/" + line
		}
		if strings.HasPrefix(line, "#") && strings.Contains(line, "URI=\"") && !strings.Contains(line, "URI=\"\"") {
			line = strings.ReplaceAll(line, "URI=\"", "URI=\""+"http://"+r.Host+"/link/"+encodedChannelName+"/")
		}
		w.Write([]byte(line + "\n"))
	}
}

func handleLinkRequest(w http.ResponseWriter, r *http.Request) {

	match := regexChannelWithURL.FindStringSubmatch(r.URL.Path)

	if match == nil {
		print404(w, "Unable to properly extract data from request '"+r.URL.Path+"'!")
		return
	}

	encodedChannelName := match[1]
	retrievedPath := match[2]

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

	resp, err := http.Get(requiredURL + retrievedPath)
	if err != nil {
		print404(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if strings.HasSuffix(r.URL.Path, ".ts") {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			print404(w, err)
			return
		}

		w.Write(body)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			line = "http://" + r.Host + "/link/" + encodedChannelName + "/" + line
		}
		if strings.Contains(line, "URI=\"") && !strings.Contains(line, "URI=\"\"") {
			line = strings.ReplaceAll(line, "URI=\"", "URI=\""+"http://"+r.Host+"/link/"+encodedChannelName+"/")
		}
		w.Write([]byte(line + "\n"))
	}

}
