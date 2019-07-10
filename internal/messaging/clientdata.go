package messaging

/* ***
 * ClientData related items
 * ClientData is used for message actions (such as CREATE)
 * and sources (such as hostname:/dir/watched/)
 * ***/

// ClientData appID for action, path for full path (including hostname)
type ClientData struct {
	appID string
	path  string
}

var clientdata ClientData

// Format appID to be in form of hostname.domains:/absolute/path/
func (cd *ClientData) get() string {
	if cd.appID != "" && cd.path != "" {
		return cd.appID + ":" + cd.path
	}

	return ""
}
