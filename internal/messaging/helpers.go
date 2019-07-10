package messaging

import "strings"

// Remove path from filename
func stripPath(filename string) string {
	return strings.Replace(filename, clientdata.path, "", 1)
}

// Fix path to always have ending forward slash
func fixPath(path string) string {
	if path[len(path)-1] != '/' {
		path = path + "/"
	}

	return path
}
