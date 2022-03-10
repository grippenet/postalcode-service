package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grippenet/postalcodes"
	"io/ioutil"
	"log"
)

func loadData(file string) (*postalcodes.PostalCodeMap, error) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	p := postalcodes.PostalCodeMap{}

	err = json.Unmarshal(data, &p)
	if err != nil {
		fmt.Println("error:", err)
	}
	return &p, nil
}

func main() {

	postalCodes, err := loadData("data/postal.json")

	if err != nil {
		log.Fatal("Unable to load postal code data. Aborting")
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/db", func(c *gin.Context) {
		c.JSON(200, gin.H{
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
