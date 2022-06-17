package main

import (
	"log"
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/grippenet/postalcodes"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}
}

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

	r.Use(Cors())

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

	r.GET("/query/:postal", func(c *gin.Context) {
		code := c.Param("postal")

		// Dont serve too long or short code
		if len(code) > 5 || len(code) < 5 {
			c.JSON(404, gin.H{"status": "Not found"})
			return
		}

		m := postalCodes.MunicipalitiesOfPostal(code)

		if m == nil {
			c.JSON(404, gin.H{"status": "Not found"})
			return
		}

		c.JSON(200, gin.H{
			"data": m,
		})

	})

	r.GET("/label/:code", func(c *gin.Context) {
		code := c.Param("code")

		if len(code) > 5 || len(code) < 5 {
			c.JSON(404, gin.H{"status": "Not found"})
			return
		}

		label := postalCodes.LabelAt(code)

		if label == "" {
			c.JSON(404, gin.H{"status": "Not found"})
			return
		}

		c.JSON(200, gin.H{
			"label": label,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
