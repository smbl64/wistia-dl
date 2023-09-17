package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type asset struct {
	URL         string  `json:"url"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Size        float64 `json:"size"`
	DisplayName string  `json:"display_name"`
}

const allowOriginal = false

var videoStreamRegex = regexp.MustCompile(`^\d+p$`)

// Start by W.iframeInit(
// ... followed by some text
// ... ended by these: , {});
var assetsRegex = regexp.MustCompile(`W\.iframeInit\((.*?)(?:,\s*{}\);)`)

func (a *asset) IsVideo() bool {
	fmt.Println(*a)
	if allowOriginal && strings.ToLower(a.DisplayName) == "original file" {
		return true
	}

	return videoStreamRegex.MatchString(a.DisplayName)
}

func findAssets(body string) ([]asset, error) {
	target := assetsRegex.FindStringSubmatch(body)
	if len(target) == 0 {
		return nil, fmt.Errorf("cannot parse the json")
	}

	jsonString := target[1]
	var output struct {
		Assets []asset
	}
	err := json.Unmarshal([]byte(jsonString), &output)
	if err != nil {
		return nil, err
	}

	return output.Assets, nil
}
