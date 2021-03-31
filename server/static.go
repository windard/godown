package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/windard/godown/fetch"
)

// StaticServerFileSystem will start HTTP Server.
func StaticServerFileSystem(host, port string, path, root string, listDirectory bool) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/", func(c *gin.Context) {
		url := c.PostForm("url")
		if url != "" {
			async := c.PostForm("async")
			if async != "" {
				go fetch.GoroutineDownload(url, 20, 10*1024*1024, 30)
			} else {
				fetch.GoroutineDownload(url, 20, 10*1024*1024, 30)
			}
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", url))
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			log.Printf("upload error:%+v", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file error:%+v", err))
			return
		}
		filename := filepath.Base(file.Filename)
		err = c.SaveUploadedFile(file, filename)
		if err != nil {
			log.Printf("upload error:%+v", err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("save file error:%+v", err))
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", filename))
	})
	router.StaticFS(path, gin.Dir(root, listDirectory))
	fmt.Printf("[GIN] Listening and serving HTTP on %s:%s\n", host, port)
	log.Fatal(router.Run(fmt.Sprintf("%s:%s", host, port)))
}
