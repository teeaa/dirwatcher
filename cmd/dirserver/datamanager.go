package main

import (
	"sort"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Directorystore contains all the data for a directory
type Directorystore struct {
	directory string
	filestore *Filestore
}

// Filestore contains data for a directory's contents
type Filestore struct {
	filedata []string
	mux      sync.Mutex
}

var directorystores []Directorystore

// Global mutex for manipulating the directory collection
var storemux sync.Mutex

// Add files to directory, creating it if it doesn't exist, and overwriting if it does
func addFilesToDir(dir string, files []string) {
	addDir(dir).addFiles(files, true)
}

func getDir(dir string) *Filestore {
	storemux.Lock()
	defer storemux.Unlock()
	for k := range directorystores {
		if directorystores[k].directory == dir {
			return directorystores[k].filestore
		}
	}

	return nil
}

func addDir(dir string) *Filestore {
	if dirStore := getDir(dir); dirStore != nil {
		return dirStore
	}
	storemux.Lock()
	directorystores = append(directorystores, Directorystore{
		directory: dir,
		filestore: &Filestore{
			filedata: []string{}, mux: sync.Mutex{},
		},
	})
	storemux.Unlock()
	return getDir(dir)
}

func rmDir(dir string) {
	storemux.Lock()
	defer storemux.Unlock()
	for k := range directorystores {
		if directorystores[k].directory == dir {
			directorystores = append(directorystores[:k], directorystores[k+1:]...)
			return
		}
	}
	log.Error("Directory ", dir, " doesn't exist")
}

func (fs *Filestore) getFiles() []string {
	return fs.filedata
}

func getFiles(dir string) []string {
	return getDir(dir).getFiles()
}

// Get a list of all the directories and their files
func getAllFiles() map[string][]string {
	filelist := make(map[string][]string, 100)
	for k := range directorystores {
		filelist[directorystores[k].directory] = getFiles(directorystores[k].directory)
	}
	return filelist
}

// Formatted version of all the files for output
func formatFiledata() []map[string]string {
	var filelist []string

	for _, files := range getAllFiles() {
		filelist = append(filelist, files...)
	}
	filelist = unique(filelist)
	sort.Strings(filelist)
	var formatted []map[string]string
	for _, file := range filelist {
		formatted = append(formatted, map[string]string{"filename": file})
	}

	// If no files we want to return an empty array instead of null for consistency
	if len(formatted) == 0 {
		formatted = []map[string]string{}
	}
	return formatted
}

func (fs *Filestore) addFile(file string) []string {
	fs.mux.Lock()
	fs.filedata = append(fs.filedata, file)
	fs.mux.Unlock()
	return fs.filedata
}

func (fs *Filestore) addFiles(files []string, replace bool) []string {
	fs.mux.Lock()
	if replace {
		fs.filedata = files
	} else {
		fs.filedata = unique(append(fs.filedata, files...))
	}
	fs.mux.Unlock()
	return fs.filedata
}

func addFile(dir string, file string) []string {
	return getDir(dir).addFile(file)
}

func (fs *Filestore) rmFile(file string) []string {
	fs.mux.Lock()
	defer fs.mux.Unlock()
	for k, v := range fs.filedata {
		if v == file {
			fs.filedata = append(fs.filedata[:k], fs.filedata[k+1:]...)
			return fs.filedata
		}
	}

	// Warn, could be just a directory being removed
	log.Warn("Unable to rm file ", file, ": File not found")

	return fs.filedata
}

func rmFile(dir string, file string) []string {
	return getDir(dir).rmFile(file)
}

// Return unique slice
func unique(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}
