package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	gp "github.com/grippenet/postalcodes"
)

// createMap create the mapping registry from the lines of the csv files
func createMap(data [][]string) *gp.PostalCodeMap {

	builder := gp.NewBuilder()

	for i, line := range data {
		if i == 0 { // omit header line
			continue
		}

		insee := line[0]
		index, found := builder.Has(insee)

		if !found {
			index = builder.Register(insee, line[1])
		}

		postal := line[2]

		builder.AddForPostal(postal, index)

	}

	return builder.GetMap()
}

func usage() {
	fmt.Println("Build postal code database and output it into json file")
	fmt.Println("Usage: build.go input_file [output_json]")
	fmt.Println("input_file: a flat csv file (semicolumn field separated) with at least 3 columns in this order municipality code;municipality label;postal code. With header but values are not checked")
	fmt.Println("output_json: Optional json output name, default is data/postal.json")
}

func main() {

	if len(os.Args) < 2 {
		usage()
		log.Fatal("File argument is required as the csv file with postal codes")

	}

	fn := os.Args[1]

	var output string = "data/postal.json"

	if len(os.Args) > 2 {
		output = os.Args[2]
	}

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
	_ = ioutil.WriteFile(output, file, 0644)
}
