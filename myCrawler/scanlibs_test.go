package main

import (
	"myCrawler/scanlibs"
	"myCrawler/utils"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestDownloadScanlibs(t *testing.T) {

	maxWorkerCount := 2
	queue := make(chan *scanlibs.ScanLibPage, maxWorkerCount)
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := sync.WaitGroup{}
	for i := 0; i < maxWorkerCount; i++ {
		go func() {
			defer wg.Done()
			wg.Add(1)
			for page := range queue {
				time.Sleep(time.Second * 1)
				setEveryPage(page)
			}
		}()
	}

	//for i := range [100]int{} {
	pages := getTotalPages()
	pages = pages[0:10]
	for _, scanLibPage := range pages {
		queue <- scanLibPage
	}
	close(queue)
	wg.Wait()

	utils.WriteJSON(pages, "scanlib_pages.json")
}
