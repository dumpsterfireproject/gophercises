package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Results struct {
	Total   int
	Correct int
}

func InitResults() *Results {
	return &Results{Total: 0, Correct: 0}
}

func (r *Results) CountQuestion() {
	r.Total = r.Total + 1
}

func (r *Results) CountCorrectAnswer() {
	r.Correct = r.Correct + 1
}

func main() {
	fileName := "problems.csv"
	flag.StringVar(&fileName, "csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.Parse()
	readCsv(fileName)
}

func readCsv(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	results := InitResults()
	csvReader := csv.NewReader(f)
	// can use csvReader.ReadAll() to return ([][]string,error)
	// can use csvReader.Comma to set different delimiter
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		results.CountQuestion()
		fmt.Printf("Problem #%d: %s = ", results.Total, rec[0])
		answer, _ := getAnswer()
		if strings.TrimSpace(answer) == rec[1] {
			results.CountCorrectAnswer()
		}
	}
	fmt.Printf("You scored %d out of %d.", results.Correct, results.Total)
}

func getAnswer() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}
