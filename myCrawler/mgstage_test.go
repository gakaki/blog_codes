package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetLinkByItemDoc(t *testing.T) {
	url := "https://sample.mgstage.com/sample/nanpatv/200gana/2396/200gana-2396_20201202T125901.ism/request?uid=10000000-0000-0000-0000-00000000000a&amp;pid=16b6ae62-e6d6-412c-afda-c8b4709c86eb"
	end := strings.Index(url, "/request?")
	lastStr := url[:end]
	lastStr = strings.Replace(lastStr, ".ism", ".mp4", 1)
	fmt.Println(lastStr)
}

func TestDownload(t *testing.T) {

	go func() {
		getList()
	}()

	go func() {
		DetailPageToImagesVideo()
	}()

	go func() {
		for video := range queueVideos {
			downloadVideo(video)
		}
	}()
	go func() {
		for video := range queueImages {
			downloadImages(video)
		}
	}()

	wg.Add(2)
	wg.Wait()

}
