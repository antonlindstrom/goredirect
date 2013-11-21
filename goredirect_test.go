package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReloadConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(ReloadConfig))
	defer server.Close()

	resp, err := http.DefaultClient.Get(server.URL + "/r/config/reload")

	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "OK" {
		t.Fatal("Didn't get OK from /r/config/reload")
	}
}

func TestStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(StatusCheck))
	defer server.Close()

	resp, err := http.DefaultClient.Get(server.URL + "/status")

	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "OK" {
		t.Fatal("Didn't get OK from /status")
	}
}

func TestHostRedirect(t *testing.T) {
	LoadConfig()
	responseCode := RunRedirect("localhost")

	if responseCode != 301 {
		t.Fatal("Returned status code was not 301, redirect failed!")
	}

	responseCode = RunRedirect("do_not_map")

	if responseCode != 503 {
		t.Fatal("Returned status code was not 503, seems like the map was found!")
	}

}

func RunRedirect(hostName string) int {
	server := httptest.NewServer(http.HandlerFunc(HostRedirect))
	defer server.Close()

	tr := &http.Transport{}

	req, err := http.NewRequest("GET", server.URL, nil)
	req.Host = hostName
	resp, err := tr.RoundTrip(req)

	if err != nil {
		log.Fatal(err)
	}

	return resp.StatusCode
}
