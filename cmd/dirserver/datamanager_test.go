package main

import (
	"strconv"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"
)

// getDir on non-existing directory should return nil
func TestGetDir(t *testing.T) {
	dir := getDir("")

	if dir != nil {
		t.Error("Expected getDir to return nil")
	}
}

// addFilesToDir should create a directory if it doesn't exist and add files to it, replacing existing data
func TestAddFilesToDir(t *testing.T) {
	dir := "dir_one"
	files := []string{"one_file_one", "one_file_two"}
	addFilesToDir(dir, files)
	gotFiles := getFiles(dir)

	if len(gotFiles) != len(files) {
		t.Error("Added files doesn't match what expected (length)")
	}

	compareFiles(t, files, gotFiles)
}

// addFile adds a single file to a directory and returns list of files in that directory
func TestAddFile(t *testing.T) {
	dir := "dir_one"
	file := "one_file_three"
	files := []string{"one_file_one", "one_file_two", file}
	added := addFile(dir, file)

	compareFiles(t, files, getFiles(dir))
	compareFiles(t, files, added)
}

// rmFile removes a single file from a directory and returns list of files in that directory
func TestRmFile(t *testing.T) {
	dir := "dir_one"
	files := []string{"one_file_one", "one_file_three"}
	file := "one_file_two"
	removed := rmFile(dir, file)

	compareFiles(t, files, getFiles(dir))
	compareFiles(t, files, removed)

	// Trying to remove the same file again should have no effect except for logging an error
	log.SetLevel(log.PanicLevel)
	removed = rmFile(dir, file)
	log.SetLevel(log.WarnLevel)
	compareFiles(t, files, removed)
}

// Test adding directory creates a directory with an initialised filestore
func TestAddDir(t *testing.T) {
	dir := "dir_two"
	addDir(dir)
	fs := getDir(dir)
	if fs == nil {
		t.Error("Expected to get added directory")
	}
	if fs.filedata == nil {
		t.Error("Expected to get empty filelist from added directory")
	}
	if len(fs.filedata) != 0 {
		t.Error("Expected filelist from added directory to be length 0")
	}
}

// Test removing a directory with and without files in it
func TestRmDir(t *testing.T) {
	dirs := []string{"dir_one", "dir_two"}
	for _, dir := range dirs {
		rmDir(dir)
		fs := getDir(dir)
		if fs != nil {
			t.Error("Expected removed dir to be nil")
		}
	}
}

// Test with load for thread safety, create 1k directories with 1k files in each
func TestMultiples(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go operations(i, &wg)
	}
	wg.Wait()
	all := getAllFiles()
	if len(all) != 1000 {
		t.Error("Expected there to be 1000 directories")
	}
	for dir, files := range all {
		if len(getDir(dir).getFiles()) != len(files) || len(files) != 1000 {
			t.Errorf("Expected there to be 1000 files in directory %s", dir)
		}
	}
}

// Remove half of the files in each directory
func TestMultiples2(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go operations2(i, &wg)
	}
	wg.Wait()
	all := getAllFiles()
	if len(all) != 1000 {
		t.Error("Expected there to be 1000 directories")
	}
	for dir, files := range all {
		if len(getDir(dir).getFiles()) != len(files) || len(files) != 500 {
			t.Errorf("Expected there to be 500 files in directory %s", dir)
		}
	}
}

// Test that adding files with replace off will only have unique results
func TestAddFiles(t *testing.T) {
	all := getAllFiles()

	for dir, files := range all {
		getDir(dir).addFiles(files, false)
		if len(getDir(dir).getFiles()) != 500 {
			t.Errorf("Expected directory to have 500 files, got %d", len(getDir(dir).getFiles()))
		}
	}
}

// Test that adding and removing files at the same time keeps the structure intact
func TestMultiples3(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go operations3(i, &wg)
	}
	wg.Wait()

	all := getAllFiles()
	for dir, files := range all {
		if len(getDir(dir).getFiles()) != len(files) || len(files) != 500 {
			t.Errorf("Expected there to be 500 files in directory %s", dir)
		}
	}
}

func operations(val int, wg *sync.WaitGroup) {
	i := strconv.Itoa(val)
	dir := "dir_" + i
	addFilesToDir(dir, []string{
		"file_1." + i,
		"file_2." + i,
		"file_3." + i,
	})
	for j := 3; j < 1000; j++ {
		addFile(dir, "file."+strconv.Itoa(j))
	}
	wg.Done()
}

func operations2(val int, wg *sync.WaitGroup) {
	i := strconv.Itoa(val)
	dir := "dir_" + i

	for j := 0; j < 500; j++ {
		rmFile(dir, getFiles(dir)[j])
	}
	wg.Done()
}

func operations3(val int, wg *sync.WaitGroup) {
	i := strconv.Itoa(val)
	dir := "dir_" + i

	files := getFiles(dir)

	for j := 0; j < 500; j++ {
		rmFile(dir, files[j])
		addFile(dir, "new_file."+strconv.Itoa(j))
	}
	wg.Done()
}

func compareFiles(t *testing.T, files []string, files2 []string) {
	for _, file := range files {
		found := false
		for _, testFile := range files2 {
			if file == testFile {
				found = true
			}
		}
		if !found {
			t.Errorf("Added files doesn't match what expected (%s not present)", file)
		}
	}
}
