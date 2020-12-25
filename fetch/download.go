package fetch

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var wait = sync.WaitGroup{}
var counter = make(map[int]int)
var counterQueue = make(chan int, 100)

// GoroutineDownload will download form requestURL.
func GoroutineDownload(requestURL string, poolSize, chunkSize int64) {
	var index, start int64

	// fetch file length
	length, err := GetFileLength(requestURL)
	if length == 0 {
		log.Printf("content not exist:%s\n", requestURL)
		return
	}

	// parse url
	u, err := url.Parse(requestURL)
	if err != nil {
		log.Printf("parse error:%+v\n", err.Error())
		return
	}
	pathList := strings.Split(u.Path, "/")
	fileName := pathList[len(pathList)-1]

	// open file
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("open error:%+v\n", err)
		return
	}
	bar := progressbar.NewOptions64(
		length,
		progressbar.OptionShowCount(),
		progressbar.OptionSetDescription(fileName),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true))

	go func() {
		for spendTime := range counterQueue {
			if n, ok := counter[spendTime]; ok {
				counter[spendTime] = n+1
			}else {
				counter[spendTime] = 1
			}
		}
	}()
	
	pool := make(chan int64, poolSize)
	for index = 0; index < poolSize; index++ {
		go func() {
			flag := true
			for flag {
				start, err := downloadChunkToFile(requestURL, pool, f, bar, chunkSize)
				if err != nil {
					log.Printf("fetch chunck start:%d error:%+v\n", start, err)
					pool <- start
				} else {
					flag = false
				}
			}
			log.Print("end one goroutine")
		}()
	}

	for start = 0; start < length; start += chunkSize {
		wait.Add(1)
		pool <- start
	}

	wait.Wait()
	fmt.Println()
	log.Printf("time spend count:%+v", counter)
}

func downloadChunkToFile(requestURL string, pool chan int64, f *os.File, bar *progressbar.ProgressBar, chunkSize int64) (start int64, err error) {
	log.Print("start new loop")
	client := &http.Client{Timeout: time.Second*30}
	chunkRequest, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Printf("create request error:%+v\n", err)
		return
	}

	for {
		startTime := time.Now()
		start := <-pool
		chunkRequest.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, start+chunkSize-1))
		resp, err := client.Do(chunkRequest)
		if err != nil {
			log.Printf("send request error:%+v\n", err)
			return start, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			_ = resp.Body.Close()
			log.Printf("read response error:%+v\n", err)
			return start, err
		}

		n, err := f.WriteAt(body, start)
		if err != nil {
			_ = resp.Body.Close()
			log.Printf("write file error:%+v\n", err)
			return start, err
		}
		spendTime := int(time.Now().Sub(startTime).Seconds())
		log.Printf("\n%d download spend %ds\n", n, spendTime)
		counterQueue <- spendTime
		_ = bar.Add(n)
		_ = resp.Body.Close()

		// echo chunk will down one.
		wait.Done()
	}
}

// GetFileLength will return http response content-length
func GetFileLength(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(resp.Status)
	}
	return resp.ContentLength, nil
}
