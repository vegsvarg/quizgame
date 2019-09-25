package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"encoding/csv"
	"strconv"
	"strings"
	"time"
)

type UserInput struct {
	problems []Problem
}

type Problem struct {
	solution int64
	userAnswer int64
	equation Equation
	correct bool
	user_correct bool
}

type Equation struct {
	first int64
	second int64
}

var problems []Problem

var FILE_PATH = "/Users/jarovit/goworkspace/src/gophercize/quizgame/data.csv"

func main()  {
	rd := openFile(FILE_PATH)
	reader := newReader(rd)
	var userAnsewrs UserInput

	var correctProblems []Problem

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		firstVal := convertToInt(strings.Split(line[0], "+")[0])
		secondVal := convertToInt(strings.Split(line[0], "+")[1])
		solution := convertToInt(line[1])

		problems = append(problems, Problem{
			solution: convertToInt(line[1]),
			equation: Equation{
				first:   firstVal,
				second:  secondVal,
			},
			correct: firstVal+secondVal == solution,
		})
	}
	correctProblems = findCorrectProblems(problems)
	userAnsewrs = timeQuiz(correctProblems, userAnsewrs)
	fmt.Println("Total number of questions answered:\t", len(userAnsewrs.problems))
	fmt.Println("Total number of questions answered correctly: \t", numberUserCorrect(userAnsewrs))
}

func convertToInt(s string) int64 {
	parsed,err := strconv.ParseInt(s,10,64)
	checkError(err)
	return parsed
}

func findCorrectProblems(p []Problem) []Problem {
	var correctProblems []Problem
	for i := 0; i < len(p) ; i++ {
		if problems[i].correct {
			correctProblems = append(correctProblems, p[i])
		}
	}
	return correctProblems
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func openFile(fileHandle string) (*os.File) {
	fileDescriptor, err := os.Open(fileHandle)
	checkError(err)
	return fileDescriptor
}

func newReader(fileDescriptor *os.File) *csv.Reader {
	reader := csv.NewReader(bufio.NewReader(fileDescriptor))
	return reader
}

func prompt()  {
	start := 0
	fmt.Println("Hello User",
		"Please answer each question one by one, you will only have thirty seconds." +
		"when you click enter the time will begin")
	fmt.Scan(&start)
	if start != 0 {
		return
	}
}


func timeQuiz(correctProblems []Problem, userInput UserInput) UserInput {
	var i int
	//var userAnswer int64

	for start := time.Now(); ; {

		if time.Since(start) > 30 * time.Second {
			return userInput
		}
		fmt.Println("Please solve: \t", correctProblems[i].equation.first, "+", correctProblems[i].equation.second, "::elapsed::\t",time.Since(start))
		fmt.Scan(&correctProblems[i].userAnswer)
		correctProblems[i].user_correct = correctProblems[i].userAnswer == correctProblems[i].solution
		userInput.problems = append(userInput.problems, correctProblems[i])

		i++
	}
}

func numberCorrect(problems []Problem) int {
	correctProblems := []Problem{}
	numCorrect := 0
	for i := 0; i < len(problems) ; i++ {
		if problems[i].correct {
			numCorrect += 1

			correctProblems = append(correctProblems, problems[i])
		}
	}
	return numCorrect
}

func numberUserCorrect(inputs UserInput) int {
	numCorrect := 0
	for i := 0; i < len(inputs.problems) ; i++ {
		if inputs.problems[i].user_correct {
			numCorrect += 1
		} else {
			fmt.Println("Equation:\t", inputs.problems[i].equation, "\nUser Answer: \t", inputs.problems[i].userAnswer, "\nActual Solution\t" ,inputs.problems[i].solution )
		}
	}
	return numCorrect
}

func printProblems(problems []Problem) {
	numCorrect := numberCorrect(problems)
	fmt.Println("Total number of problems: \t", len(problems))
	fmt.Println("Total number correct: \t", numCorrect)
}
