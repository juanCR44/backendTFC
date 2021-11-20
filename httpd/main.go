package main

import (
	"Backend/httpd/handler"
	"Backend/platform/newsfeed"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	feed := newsfeed.New()
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/ping", handler.RouteGet())
	router.GET("/newsfeed", handler.NewsfeedGet(feed))
	router.POST("/newsfeed", handler.NewsfeedPost(feed))

	router.POST("/upload", func(c *gin.Context) {
		fmt.Println(c.ContentType())
		file, err := c.FormFile("file")

		//fmt.Println(file, " aver ", err)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		filename := filepath.Base(file.Filename)
		fmt.Printf("Uploaded File : %+v\n", filename)

		out, _ := os.Open(file.Filename)

		csvLines, readErr := csv.NewReader(out).ReadAll()
		if readErr != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		for _, line := range csvLines {
			fmt.Println(line)
		}
	})

	router.Run("localhost:8080")
}
