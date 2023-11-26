package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("tcp", "1.1.1.1:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.TCPAddr)

	return localAddr.IP
}

func headersByRequestURI() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/dl") {
			c.Header("Content-Disposition", "attachment")
		}
	}
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Use(headersByRequestURI())
	router.StaticFS("/dl", http.Dir("./dl"))
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/list")
	})
	router.GET("/list/*action", func(c *gin.Context) {
		action := c.Param("action")
		actionDir, err := os.ReadDir("./dl/" + action)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		}
		var files []string
		var dirs []string
		for _, ele := range actionDir {
			if ele.IsDir() {
				dirs = append(dirs, ele.Name())
			} else {
				files = append(files, ele.Name())
			}
		}
		var path string
		if action == "/" {
			path = ""
		} else {
			path = action
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":  "临时下载站",
			"action": path,
			"files":  files,
			"dirs":   dirs,
		})
	})

	println(GetOutboundIP().String())

	router.Run(":28100")
}
