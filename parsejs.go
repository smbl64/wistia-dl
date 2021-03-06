package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type asset struct {
	url           string
	width, height int
	size          float64
	displayName   string

	isVideo bool
}

func isVideoStream(a *asset, allowOriginal bool) bool {
	if allowOriginal && strings.ToLower(a.displayName) == "original file" {
		return true
	}

	reg := regexp.MustCompile(`^\d+p$`)
	return reg.MatchString(a.displayName)
}

func findAssets(body string) ([]asset, error) {
	// Start by W.iframeInit(
	// ... followed by some text
	// ... ended by these: , {});
	reg := regexp.MustCompile(`W\.iframeInit\((.*?)(?:,\s*{}\);)`)

	target := reg.FindStringSubmatch(body)
	if len(target) == 0 {
		return nil, fmt.Errorf("cannot parse the json")
	}

	jsonString := target[1]
	var parsedJson map[string]interface{}
	json.Unmarshal([]byte(jsonString), &parsedJson)

	if _, ok := parsedJson["assets"]; !ok {
		return nil, fmt.Errorf("cannot find assets in the json")
	}

	jsonAssets := parsedJson["assets"].([]interface{})

	assets := make([]asset, 0)
	for _, value := range jsonAssets {
		row := value.(map[string]interface{})

		a := parseAssetRow(row)
		assets = append(assets, a)
	}
	return assets, nil
}

func parseAssetRow(row map[string]interface{}) asset {
	var a asset

	hasField := func(key string) bool {
		_, ok := row[key]
		return ok
	}

	a.url = row["url"].(string)
	a.displayName = row["display_name"].(string)
	if hasField("width") {
		a.width = int(row["width"].(float64))
	}
	if hasField("height") {
		a.height = int(row["height"].(float64))
	}
	if hasField("size") {
		a.size = row["size"].(float64)
	}
	a.isVideo = isVideoStream(&a, false)

	return a
}
