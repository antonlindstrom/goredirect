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

	http.HandleFunc("/r/config/dump", DumpConfig)
	http.HandleFunc("/r/config/reload", ReloadConfig)
	http.HandleFunc("/r/status", StatusCheck)
	http.HandleFunc("/", HostRedirect)

	logger("Serving requests on port %s\n", port)
	log.Fatal(s.ListenAndServe())
}

// Redirect to map
func HostRedirect(w http.ResponseWriter, req *http.Request) {
	requestUrl := strings.Split(req.Host, ":")[0]
	redirectUrl := RedirectMap[requestUrl]

	if redirectUrl == "" {
		logger("code: 503, request: %s, redirect: empty!, path: %s", requestUrl, req.URL.Path)
		http.Error(w, "503: Could not map request!", 503)
		return
	}

	logger("code: %d, request: %s, redirect: %s, path: %s", RedirectType, requestUrl, redirectUrl, req.URL.Path)
	http.Redirect(w, req, redirectUrl+req.URL.Path, RedirectType)
}

// Dump RedirectMap
func DumpConfig(w http.ResponseWriter, req *http.Request) {
	json, _ := json.Marshal(RedirectMap)
	logger("code: 200, dump")
	fmt.Fprintf(w, string(json))
}

// Handler to reload configuration
func ReloadConfig(w http.ResponseWriter, req *http.Request) {
	msg := LoadConfig()
	logger("code: 200, reload")
	fmt.Fprintf(w, msg)
}

// Status handler
func StatusCheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK")
}

// Load configuration
func LoadConfig() string {
	bytes, err := ioutil.ReadFile("./config.json")

	if err != nil {
		log.Printf("Error %s\n", err)
		return "Could not read config file"
	}

	err = json.Unmarshal(bytes, &RedirectMap)

	if err != nil {
		log.Printf("Error %s\n", err)
		return "JSON parse error"
	}

	return "OK"
}

// Log if we're running verbose
func logger(message string, args ...interface{}) {
	if Verbose {
		log.Printf(message, args...)
	}
}
