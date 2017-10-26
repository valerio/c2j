package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "c2j"
	app.Description = "A simple cli tool for converting csv to json"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Input file to be converted",
		},
	}

	app.Action = func(c *cli.Context) error {
		run(c.String("file"))
		return nil
	}

	app.Run(os.Args)
}

func run(inputFile string) {
	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(bufio.NewReader(f))

	headers, err := csvReader.Read()
	if err != nil {
		panic(err)
	}

	result, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	encoder := json.NewEncoder(os.Stdout)

	fmt.Println("[")
	first := true

	for _, parsed := range result {
		m := make(map[string]string)

		for idx, header := range headers {
			if parsed[idx] == "" {
				continue
			}

			m[header] = parsed[idx]
		}

		if first {
			first = !first
		} else {
			fmt.Print(",")
		}

		encoder.Encode(m)
	}

	fmt.Print("]")
}
