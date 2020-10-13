package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	choices  []string
	answer   string
}

// Creates flags for the program.
func parseFlag() (*string, *int) {
	csvFileName := flag.String("csv", "quiz.csv", "CSV File for the quiz. Default: quiz.csv")
	duration := flag.Int("duration", 120, "Duration for the whole quiz in seconds.  Default: 120 seconds.")
	flag.Parse()
	return csvFileName, duration
}

// Exit the program with the specified message.
func exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

// Read each of the files' lines. Returns the lines.
func readFilesLine(file *os.File) [][]string {
	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit(fmt.Sprintf("ERROR! Couldn't read file! %s", err.Error()))
	}

	return lines
}

// Parses the lines into the question, choices, and answers.
func parseLines(lines [][]string) []problem {
	retVal := make([]problem, len(lines))
	for i, line := range lines {
		retVal[i] = problem{
			question: line[0],
			choices:  []string{line[1], line[2], line[3], line[4]},
			answer:   strings.TrimSpace(line[5]),
		}
	}
	return retVal
}

// Invoke/Run the quiz. Returns score and the quantity of the quiz.
func invokeQuiz(problems []problem, timer *time.Timer) (int, int) {
	score := 0
	s := bufio.NewScanner(os.Stdin)

	for i, p := range problems {
		fmt.Printf("%d. %s \n%s\n%s\n%s\n%s\n", i+1, p.question, p.choices[0],
			p.choices[1], p.choices[2], p.choices[3])
		fmt.Printf("Your answer: ")
		answerCh := make(chan string)

		go func() {
			s.Scan()
			answer := s.Text()
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\n\nYour time is up.")
			return score, len(problems)
		case answer := <-answerCh:
			if strings.ToUpper(answer) == p.answer {
				score++
			}
			fmt.Println()
		}
	}
	return score, len(problems)
}

// Show the user's score.
func showScore(score int, quantity int) {
	s := bufio.NewScanner(os.Stdin)

	fmt.Printf("Calculating score...")

	time.Sleep(1 * time.Second)
	fmt.Printf("\nYou scored %d out of %d.", score, quantity)

	fmt.Println("\n\nPlease press Enter to exit...")
	s.Scan()
}

func main() {
	csvFile, duration := parseFlag()
	file, err := os.Open(*csvFile)

	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV File: %s", *csvFile))
	}

	lines := readFilesLine(file)
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*duration) * time.Second)
	score, qty := invokeQuiz(problems, timer)

	showScore(score, qty)
}
