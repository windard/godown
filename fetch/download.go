/*

Package fetch for download, provide high performance download

use Goroutine to parallel download, use WaitGroup to do concurrency control.

*/
package fetch

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
)

// FileFlag save file flag
//
// FileMode save file mode
const (
	FileFlag = os.O_WRONLY | os.O_CREATE
	FileMode = 0644
)

// WaitPool implement request pool to enhance performance
var (
	WaitPool = sync.WaitGroup{}
)

// GoroutineDownload will download form requestURL.
// example:
//  requestURL := "http://xxx"
//  GoroutineDownload(requestURL, 20, 10*1024*1024, 30)
func GoroutineDownload(requestURL string, poolSize, chunkSize, timeout int64, force bool) {
	var index, start int64

	if !strings.HasPrefix(requestURL, "http") {
		requestURL = "http://" + requestURL
	}
	requestURL = strings.TrimSpace(requestURL)

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

	if _, err = os.Stat(fileName); !os.IsNotExist(err) {
		if force {
			err = os.Remove(fileName)
			if err != nil {
				return
			}
		} else {
			fmt.Printf("[%s] is already downloaded.\n", fileName)
			return
		}
	}

	// open file
	f, err := os.OpenFile(fileName, FileFlag, FileMode)
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

	pool := make(chan int64, (length/chunkSize)+1)
	for index = 0; index < poolSize; index++ {
		go func() {
			// recover
			defer func() {
				if err2 := recover(); err2 != nil {
					log.Printf("panic error: %+v, stack:%s", err2, debug.Stack())
				}
			}()

			// loop download until finish
			for {
				start, err = downloadChunkToFile(requestURL, pool, f, bar, chunkSize, timeout)
				if err != nil {
					log.Printf("fetch chunck start:%d error:%+v\n", start, err)
					pool <- start
				} else {
					break
				}
				log.Printf("start loop download again")
			}
		}()
	}

	for start = 0; start < length; start += chunkSize {
		WaitPool.Add(1)
		pool <- start
	}

	WaitPool.Wait()
	fmt.Println()
}

func downloadChunkToFile(requestURL string, pool chan int64, f *os.File, bar *progressbar.ProgressBar, chunkSize, timeout int64) (start int64, err error) {
	client := &http.Client{Timeout: time.Second * time.Duration(timeout)}
	chunkRequest, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Printf("create request error:%+v\n", err)
		return
	}

	var resp *http.Response
	var body []byte
	var written int
	for {
		start = <-pool
		chunkRequest.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, start+chunkSize-1))
		resp, err = client.Do(chunkRequest)
		if err != nil {
			log.Printf("send request error:%+v\n", err)
			return
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			_ = resp.Body.Close()
			log.Printf("read response error:%+v\n", err)
			return
		}

		written, err = f.WriteAt(body, start)
		if err != nil {
			_ = resp.Body.Close()
			log.Printf("write file error:%+v\n", err)
			return
		}
		_ = bar.Add(written)
		_ = resp.Body.Close()

		// echo chunk will down one.
		WaitPool.Done()
	}
}

// GetFileLength will return http response content-length
// example:
//  GetFileLength("http://xxx")
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
