package release

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitlab.com/brandonbutler/chiabot/internal/helpers"
	"golang.org/x/net/html"
)

var (
	latestAPI = "https://latest.cmm.io/chia"
)

//GetLatest queries an API to receive the latest release
func GetLatest(client *http.Client) (string, error) {
	resp, stat, err := helpers.HTTPRequest(client, latestAPI, "GET", nil)
	if err != nil {
		return "", fmt.Errorf("[ERROR] GetLatest: %v", err)
	}
	if stat != 200 {
		return "", fmt.Errorf("[ERROR] GetLatest: status code returned %d", stat)
	}

	return string(resp), nil
}

type Changelog struct {
	Added   []string
	Changed []string
	Fixed   []string
}

func GetChanges(client *http.Client, url string) (Changelog, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Changelog{}, fmt.Errorf("[ERROR] GetChanges: %v", err)
	}
	defer resp.Body.Close()

	cl := Changelog{}
	//recordingChanges while crawling html nodes, this looks for supported Changelog headers. Supports "Added", "Changed", and "Fixed"
	var recordingChanges string
	sameBullet := false

	htmlTokens := html.NewTokenizer(resp.Body)
	for {
		tt := htmlTokens.Next()
		t := htmlTokens.Token()

		err := htmlTokens.Err()
		if err == io.EOF {
			break
		}

		switch tt {
		case html.StartTagToken, html.EndTagToken:
			if recordingChanges != "" {
				if tt.String() == "EndTag" && t.String() == "</ul>" {
					recordingChanges = ""
				}
				if t.String() == "</div>" {
					recordingChanges = ""
				}
				if t.String() == "</li>" {
					sameBullet = false
				}

			}
		case html.TextToken:
			data := strings.TrimSpace(t.Data)
			if strings.EqualFold("Added", data) || strings.EqualFold("Changed", data) || strings.EqualFold("Fixed", data) {
				recordingChanges = data
			} else {
				if strings.EqualFold("Added", recordingChanges) {
					if sameBullet && len(cl.Added) > 1 {
						lastElementIndex := len(cl.Added) - 1
						cl.Added[lastElementIndex] += fmt.Sprintf(" %s", data)
					} else {
						cl.Added = append(cl.Added, data)
						sameBullet = true
					}
				}
				if strings.EqualFold("Changed", recordingChanges) {
					if sameBullet && len(cl.Changed) > 1 {
						lastElementIndex := len(cl.Changed) - 1
						cl.Changed[lastElementIndex] += fmt.Sprintf(" %s", data)
					} else {
						cl.Changed = append(cl.Changed, data)
						sameBullet = true
					}
				}
				if strings.EqualFold("Fixed", recordingChanges) {
					if sameBullet && len(cl.Fixed) > 1 {
						lastElementIndex := len(cl.Fixed) - 1
						cl.Fixed[lastElementIndex] += fmt.Sprintf(" %s", data)
					} else {
						cl.Fixed = append(cl.Fixed, data)
						sameBullet = true
					}
				}
			}
		}
	}
	return cl, nil
}
