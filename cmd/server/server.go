package main

import (
	"log"
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/grippenet/postalcodes"
)

func main() {

	data_file := os.Getenv("POSTAL_FILE")
	if data_file == "" {
		data_file = "data/postal.json"
	}

	log.Printf("Using file '%s'", data_file)

	postalCodes, err := postalcodes.LoadPostalCodeMap(data_file)

	if err != nil {
		log.Fatal("Unable to load postal code data. Aborting")
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/db", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"postal_codes": gin.H{
				"codes":          len(postalCodes.Postalcodes),
				"municipalities": len(postalCodes.Municipalities),
			},
		})
	})

	r.GET("/postal/:code", func(c *gin.Context) {
		code := c.Param("code")

		m := postalCodes.MunicipalitiesOf(code)

		if m == nil {
			c.JSON(404, gin.H{"status": "Not found"})
			return
		}

		c.JSON(200, gin.H{
			"data": m,
		})

	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
