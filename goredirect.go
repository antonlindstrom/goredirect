package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const RedirectType = 301

var Verbose bool
var RedirectMap map[string]string

// Server, make sure we fail as fast as possible
func main() {
	var port = "8080"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	flag.BoolVar(&Verbose, "verbose", true, "Set verbose output")
	flag.Parse()

	LoadConfig()

	s := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/", HostRedirect)
	http.HandleFunc("/reload", ReloadConfig)
	http.HandleFunc("/status", StatusCheck)

	if Verbose {
		log.Printf("Serving requests on port %s\n", port)
	}

	log.Fatal(s.ListenAndServe())
}

// Redirect to map
func HostRedirect(w http.ResponseWriter, req *http.Request) {
	requestUrl := strings.Split(req.Host, ":")[0]
	redirectUrl := RedirectMap[requestUrl]

	if redirectUrl == "" {
		if Verbose {
			log.Printf("code: 503, request: %s, redirect: empty!, path: %s", requestUrl, req.URL.Path)
		}
		http.Error(w, "503: Could not map request!", 503)
		return
	}

	if Verbose {
		log.Printf("code: %d, request: %s, redirect: %s, path: %s", RedirectType, requestUrl, redirectUrl, req.URL.Path)
	}
	http.Redirect(w, req, redirectUrl, RedirectType)
}

// Handler to reload configuration
func ReloadConfig(w http.ResponseWriter, req *http.Request) {
	LoadConfig()
	if Verbose {
		log.Printf("code: 200, reload")
	}
	fmt.Fprintf(w, "OK, reloaded.")
}

// Status handler
func StatusCheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK")
}

// Load configuration
func LoadConfig() {
	bytes, err := ioutil.ReadFile("./config.json")

	if err != nil {
		log.Fatal("Error ", err)
	}

	err = json.Unmarshal(bytes, &RedirectMap)

	if err != nil {
		log.Fatal("Error ", err)
	}
}
