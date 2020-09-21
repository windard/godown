package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)


func StaticServerFileSystem(host, port string, path, root string, listDirectory bool) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.StaticFS(path, gin.Dir(root, listDirectory))
	fmt.Printf("[GIN] Listening and serving HTTP on %s:%s\n", host, port)
	log.Fatal(router.Run(fmt.Sprintf("%s:%s", host, port)))
}
