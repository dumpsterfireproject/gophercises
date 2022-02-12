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
	"time"
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
	var fileName string
	var timeLimit int
	flag.StringVar(&fileName, "csv", "problems.csv", "a csv file in the format of 'question,answer' (default problems.csv)")
	flag.IntVar(&timeLimit, "limit", 30, "the time limit for the quiz in seconds (default 30)")
	flag.Parse()
	readCsv(fileName, timeLimit)
}

func readCsv(fileName string, timeLimit int) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	results := InitResults()
	csvReader := csv.NewReader(f)
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	fmt.Println("Press enter to start quiz")
	getAnswer()
	go func() {
		<-timer.C
		printResults(results)
		os.Exit(0)
	}()
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
	printResults(results)
}

func getAnswer() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}

func printResults(results *Results) {
	fmt.Printf("You scored %d out of %d.\n", results.Correct, results.Total)
}
