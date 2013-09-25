package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const RedirectType = 301

var RedirectMap map[string]string

// Server, make sure we fail as fast as possible
func main() {
	var port = "8080"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	loadConfiguration()

	s := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/", HostRedirect)
	http.HandleFunc("/reload", ReloadConfig)

	log.Printf("Serving requests on port %s\n", port)
	log.Fatal(s.ListenAndServe())
}

// Redirect to map
func HostRedirect(w http.ResponseWriter, req *http.Request) {
	requestUrl := strings.Split(req.Host, ":")[0]
	redirectUrl := RedirectMap[requestUrl]

	if redirectUrl == "" {
		http.Error(w, "503: Could not map request!", 503)
		log.Printf("code: 503, request: %s, redirect: empty!, path: %s", requestUrl, req.URL.Path)
		return
	}

	log.Printf("code: %d, request: %s, redirect: %s, path: %s", RedirectType, requestUrl, redirectUrl, req.URL.Path)
	http.Redirect(w, req, redirectUrl, RedirectType)
}

// Handler to reload configuration
func ReloadConfig(w http.ResponseWriter, req *http.Request) {
	loadConfiguration()
	log.Printf("code: 200, reload")
	fmt.Fprintf(w, "OK, reloaded.")
}

// Load configuration
func loadConfiguration() {
	bytes, err := ioutil.ReadFile("./config.json")

	if err != nil {
		log.Printf("%s\n", err)
		os.Exit(1)
	}

	parseErr := json.Unmarshal(bytes, &RedirectMap)

	if parseErr != nil {
		log.Fatal("%s\n", err)
		os.Exit(1)
	}
}
