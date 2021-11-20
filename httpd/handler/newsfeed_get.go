package handler

import (
	"Backend/platform/newsfeed"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewsfeedGet(feed newsfeed.Getter) gin.HandlerFunc {

	return func(c *gin.Context) {
		/*c.JSON(http.StatusOK, map[string]string{
			"hello": "found me",
		})*/
		results := feed.GetAll()
		c.JSON(http.StatusOK, results)
	}
}
