package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format 'question,answer'")
	duration := flag.Int("duration", 30, "The quiz duration before it expires")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV File with name %s \n", *csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse provided CSV")
	}

	// Parse lines and get them as a problem struct
	problems := parseLines(lines)

	// Set up a timer duration in seconds
	timer := time.NewTicker(time.Duration(*duration) * time.Second)

	correct := 0

	// Loop through our list of problems
	for i, p := range problems {

		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)

		// set up an answer channel
		answerCh := make(chan string)

		// Run a go routine that will ask the question and wait for the answer
		go func() {
			var answer string
			_, e := fmt.Scanf("%s\n", &answer)
			if e != nil {
				fmt.Println("Error reading user entry")
			}
			// send the answer to the answer channel
			answerCh <- answer
		}()

		// Using a select to wait for messages from the channels
		select {

		// If we have a message from the Timer Channel
		case <-timer.C:
			// Show message and stop execution
			fmt.Println("Duration has expired!")
			return

		// If we have a message from the AnswerChannel, check if it is correct and update the counters
		case answer := <-answerCh:

			if answer == p.answer {
				correct++
			}
		}
	}

	fmt.Printf("You got %d answers correct out of %d\nw", correct, len(lines))
}

type problem struct {
	question string
	answer   string
}

// Parser the two dimensional array of lines, and return an array of problem structs
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: strings.TrimSpace(line[0]),
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
