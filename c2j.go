package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "c2j"
	app.UsageText = "./c2j input.csv [options] > output.json"
	app.Description = "A simple cli tool for converting csv to json"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "noheaders, nh",
			Usage: "specifies that the input file has no headers",
		},
	}

	app.Action = func(c *cli.Context) error {
		fileName := c.Args().First()
		if fileName == "" {
			return errors.New("Must specify an input file")
		}

		return run(fileName, c.Bool("noheaders"))
	}

	err := app.Run(os.Args)

	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
}

func run(inputFile string, noHeaders bool) error {
	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	defer f.Close()

	csvReader := csv.NewReader(bufio.NewReader(f))
	var headers []string

	if noHeaders == false {
		headers, err = csvReader.Read()
		if err != nil {
			return err
		}
	}

	result, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)

	writer.WriteString("[\n")
	first := true

	for _, parsed := range result {
		if first {
			first = !first
		} else {
			writer.WriteString(",")
		}

		m := make(map[string]string)
		if headers != nil {
			for idx, header := range headers {
				if parsed[idx] == "" {
					delete(m, header)
					continue
				}

				m[header] = parsed[idx]
			}
			encoder.Encode(m)
		} else {
			encoder.Encode(parsed)
		}

	}

	writer.WriteString("]")
	return nil
}
