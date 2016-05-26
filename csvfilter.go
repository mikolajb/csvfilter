package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type filters []filter

type filter struct {
	column int
	value  string
}

func (f *filters) Set(filterspec string) error {
	x := strings.Split(filterspec, "=")
	if len(x) != 2 {
		return errors.New("Wrong format")
	}
	column, err := strconv.ParseInt(x[0], 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	if column < 0 {
		return errors.New("Column number has to be >= 0")
	}
	value := x[1]
	*f = append(*f, filter{int(column), value})
	return nil
}

func (f *filters) String() string {
	result := ""
	for _, filter := range *f {
		result += fmt.Sprintf("%d=%s\n", filter.column, filter.value)
	}
	return result
}

type output []int

func (o *output) Set(outputspec string) error {
	x := strings.Split(outputspec, ",")
	for _, i := range x {
		column, err := strconv.ParseInt(i, 10, 0)
		if err != nil {
			log.Fatal(err)
		}
		if column < 0 {
			return errors.New("Column number has to be >= 0")
		}
		*o = append(*o, int(column))
	}
	return nil
}

func (o *output) String() string {
	return fmt.Sprint(*o)
}

func prepareOutput(record []string, fields []int, separator rune) string {
	output := record
	if len(fields) > 0 {
		output = make([]string, 0)
		for _, i := range fields {
			output = append(output, record[i])
		}
	}
	for i := range output {
		if strings.Contains(output[i], ",") {
			output[i] = "\"" + output[i] + "\""
		}
	}
	return strings.Join(output,
		fmt.Sprintf("%c", separator))
}

func main() {
	var filtersFlag filters
	var outputFlag output
	flag.Var(&filtersFlag, "filter", "column_num=value")
	flag.Var(&outputFlag, "select", "column_num,")
	enumheader := flag.Bool("enumheader", false, "-enumheader")
	flag.Parse()
	reader := csv.NewReader(bufio.NewReader(os.Stdin))
	header := true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if *enumheader {
			for i, value := range record {
				fmt.Printf("%3d: %s\n", i, value)
			}
			break
		}

		allOk := true
		for _, filter := range filtersFlag {
			if len(record) < filter.column {
				allOk = false
				log.Fatal("Not enough columns")
				break
			}
			if record[filter.column] != filter.value {
				allOk = false
				break
			}
		}
		if allOk || header {
			fmt.Println(prepareOutput(record, outputFlag, reader.Comma))
			header = false
		}
	}
}
