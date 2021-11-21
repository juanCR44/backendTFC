package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/upload", func(c *gin.Context) {
		fmt.Println(c.ContentType())

		file, _, err := c.Request.FormFile("file")

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		csvlines, err := csv.NewReader(file).ReadAll()

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		fmt.Print(csvlines)

	})

	bufferIn := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto remoto: ")
	puerto, _ := bufferIn.ReadString('\n')
	puerto = strings.TrimSpace(puerto)
	remotehost := fmt.Sprintf("localhost:%s", puerto)
	con, _ := net.Dial("tcp", remotehost)
	defer con.Close()
	fmt.Fprintln(con, 35)

	router.Run("localhost:8080")
}
