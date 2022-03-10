package main

import (
	"encoding/csv"
	"encoding/json"
	gp "github.com/grippenet/postalcodes"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func createMap(data [][]string) gp.PostalCodeMap {

	municipalities := make(map[string]string)
	postalcodes := make(map[string][]string)

	for i, line := range data {
		if i == 0 { // omit header line
			continue
		}

		insee := line[0]

		_, found := municipalities[insee]

		if !found {
			municipalities[insee] = line[1]
		}

		postal := line[2]

		var pp []string
		var ok bool

		pp, ok = postalcodes[postal]

		if !ok {
			pp = make([]string, 0)
		}

		pp = append(pp, insee)

		postalcodes[postal] = pp
	}

	return gp.PostalCodeMap{
		BuiltAt:        time.Now(),
		Postalcodes:    postalcodes,
		Municipalities: municipalities,
	}
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("File argument is required as the csv file with postal codes")
	}

	fn := os.Args[1]

	log.Printf("Reading file %s", fn)
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	postalMap := createMap(data)

	file, _ := json.Marshal(postalMap)
	_ = ioutil.WriteFile("data/postal.json", file, 0644)
}
