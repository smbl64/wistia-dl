package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/jessevdk/go-flags"
)

const userAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36"

var opts struct {
	VideoID  string `short:"v" long:"video-id" description:"Video ID" required:"true"`
	FileName string `short:"o" long:"output" description:"Output file" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	err = downloadVideo(opts.VideoID, opts.FileName)
	if err != nil {
		fmt.Printf("Failed: %s", err)
		os.Exit(1)
	}
}

func downloadVideo(videoID, filename string) (err error) {
	url := fmt.Sprintf("https://fast.wistia.net/embed/iframe/%s?videoFoam=true", videoID)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	bodyString := string(bodyBytes)

	assets, err := findAssets(bodyString)
	if err != nil {
		return fmt.Errorf("cannot find the target URL. Bad video ID maybe?")
	}

	asset, err := chooseAsset(assets)
	if err != nil {
		return err
	}

	fmt.Printf("Video found. Resolution=%dx%d  Size=%s\n", asset.width, asset.height, humanize.Bytes(uint64(asset.size)))

	return downloadFile(asset.url, filename)
}

// chooseAsset finds a video stream with the highest resolution
func chooseAsset(assets []asset) (asset, error) {
	var chosen asset

	for _, a := range assets {
		if !a.isVideo {
			continue
		}

		if a.height > chosen.height {
			chosen = a
		}
	}

	if chosen.height > 0 {
		return chosen, nil
	}

	return chosen, fmt.Errorf("there is no video stream in the assets")
}

func downloadFile(url, filename string) (err error) {
	fmt.Println("Starting download...")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded %d bytes.\n", n)

	return nil
}
