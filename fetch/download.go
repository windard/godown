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
)

var wait = sync.WaitGroup{}


func GoroutineDownload(requestUrl string, poolSize, chunkSize int64) {
	var start int64 = 0

	// fetch file length
	length, err := GetFileLength(requestUrl)
	if length == 0 {
		log.Printf("content not exist:%s\n", requestUrl)
		return
	}

	// parse url
	u, err := url.Parse(requestUrl)
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

	pool := make(chan int64, poolSize)
	for start = 0; start < poolSize; start++ {
		go func() {
			start,e := DownloadChunkToFile(requestUrl, pool, f, bar, chunkSize)
			log.Printf("fetch chunck start:%d error:%+v\n", start, e)
			wait.Add(1)
			pool <- start
		}()
	}

	for start = 0; start < length; start+=chunkSize {
		wait.Add(1)
		pool <- start
	}

	wait.Wait()
}

func DownloadChunkToFile(requestUrl string, pool chan int64, f *os.File, bar *progressbar.ProgressBar, chunkSize int64) (start int64, err error) {
	client := &http.Client{}
	chunkRequest, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Printf("create request error:%+v\n", err)
		return
	}

	for {
		start := <- pool
		chunkRequest.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, start+chunkSize-1))
		resp, err := client.Do(chunkRequest)
		if err != nil {
			log.Printf("send request error:%+v\n", err)
			return start, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("read response error:%+v\n", err)
			return start, err
		}

		n, err := f.WriteAt(body, start)
		if err != nil {
			log.Printf("write file error:%+v\n", err)
			return start, err
		}

		err = bar.Add(n)
		_ = resp.Body.Close()
		wait.Done()
	}
}

func GetFileLength(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		log.Println(err)
		return 0, err
	}else {
		if resp.StatusCode != http.StatusOK {
			return 0, errors.New(resp.Status)
		}
		return resp.ContentLength, nil
	}
}