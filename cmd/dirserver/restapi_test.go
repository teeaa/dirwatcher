package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var srv *http.Server
var files map[string][]map[string]string

func TestNonExisting(t *testing.T) {
	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Error("Error making request:", err)
	}
	responseReq := httptest.NewRecorder()
	router := getRouter()
	router.ServeHTTP(responseReq, req)

	if responseReq.Code != http.StatusNotFound {
		t.Errorf("Expected HTTP code 404, got %d", responseReq.Code)
	}
}

func TestEmptyFiles(t *testing.T) {
	directorystores = []Directorystore{} // reset datastore

	err := json.Unmarshal(HTTPGetFiles(t), &files)
	if err != nil {
		t.Error("Failed parsing json from HTTP server")
	}

	if len(files["files"]) != 0 {
		t.Errorf("Expected response to have empty array of files, but got %d: %s", len(files["files"]), files["files"])
	}
}

func TestFilesWithData(t *testing.T) {
	addFilesToDir("testdir", []string{"dog.jpg", "cat.gif", "hedgehog.png"})

	err := json.Unmarshal(HTTPGetFiles(t), &files)
	if err != nil {
		t.Error("Failed parsing json from HTTP server")
	}

	if len(files["files"]) != 3 {
		t.Errorf("Expected response to have three files, got %d: %s", len(files["files"]), files["files"])
	}

	filelist := files["files"]
	if filelist[0]["filename"] != "cat.gif" || filelist[1]["filename"] != "dog.jpg" || filelist[2]["filename"] != "hedgehog.png" {
		t.Error("Expected files to be sorted alphabetically but got", filelist)
	}
}

func TestFilesWithDuplicates(t *testing.T) {
	addFilesToDir("testdir2", []string{"hedgehog.png"})

	err := json.Unmarshal(HTTPGetFiles(t), &files)
	if err != nil {
		t.Error("Failed parsing json from HTTP server")
	}

	if len(files["files"]) != 3 {
		t.Errorf("Expected response to have three files, got %d: %s", len(files["files"]), files["files"])
	}
}

// Gets /files, checks response code and returns body
func HTTPGetFiles(t *testing.T) []byte {
	req, err := http.NewRequest("GET", "/files", nil)
	if err != nil {
		t.Error("Error making request:", err)
	}
	responseReq := httptest.NewRecorder()
	router := getRouter()
	router.ServeHTTP(responseReq, req)

	if responseReq.Code != http.StatusOK {
		t.Errorf("Expected HTTP code 200, got %d", responseReq.Code)
	}

	return responseReq.Body.Bytes()
}
