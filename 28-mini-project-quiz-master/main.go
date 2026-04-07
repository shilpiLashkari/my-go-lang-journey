package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		// If default file doesn't exist, create a sample one for the user
		if *csvFilename == "problems.csv" {
			createSampleCSV(*csvFilename)
			file, _ = os.Open(*csvFilename)
		} else {
			exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
		}
	}
	defer file.Close()

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	fmt.Printf("🧠 Welcome to Quiz Master!\n")
	fmt.Printf("You have %d seconds to answer %d questions.\n", *timeLimit, len(problems))
	fmt.Println("Press Enter to start!")
	fmt.Scanln()

	correct := 0

problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\n\n⏰ TIME'S UP!")
			break problemLoop
		case answer := <-answerCh:
			if strings.EqualFold(strings.TrimSpace(answer), p.a) {
				correct++
			}
		}
	}

	fmt.Printf("\n🏁 QUIZ COMPLETE!\n")
	fmt.Printf("Score: %d out of %d\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func createSampleCSV(filename string) {
	content := "5+5,10\n7*2,14\n12-8,4\nGo Creator (Rob/Ken/Robert),Rob\nWhat is Go case?,C\n"
	os.WriteFile(filename, []byte(content), 0644)
	fmt.Printf("📝 Created sample quiz file: %s\n", filename)
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
