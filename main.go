package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

func problemPuller(filename string) ([]problem, error) {
	// read all the problems from the quiz.csv file
	// 1. open the file
	if fObj, err := os.Open(filename); err == nil {
		// 2. we will create a new reader
		csvR := csv.NewReader(fObj)
		// 3. it will read the file
		if cLines, err := csvR.ReadAll(); err == nil {
			// 4. call the parse problem function
			return parseProblem(cLines), nil
		} else {

			return nil, fmt.Errorf("error in reading data in csv"+"format from %s file: %s", filename, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening %s file: %s", filename, err.Error())
	}
}

func parseProblem(lines [][]string) []problem {

	// go over the lines and parse them, with problem struct
	r := make([]problem, len(lines))

	for i := 0; i < len(lines); i++ {
		r[i] = problem{
			question: lines[i][0],
			answer:   lines[i][1],
		}
	}

	return r
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	// 1. input the name of the file
	fName := flag.String("f", "quiz.csv", "path of the csv file")

	// 2. Set the duration for the timer
	timer := flag.Int("t", 30, "time for the quiz")
	flag.Parse()
	// 3. Pull the problems from the file (calling our problem puller function)
	problems, err := problemPuller(*fName)
	// 4. handle the error
	if err != nil {
		exit(fmt.Sprintf("something went wrong:%s", err.Error()))
	}
	// 5. create a variable to count our correct answer
	correctAnswer := 0
	// 6.using the duration of the timer, we want to initialize the timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)

	ansC := make(chan string)

problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, p.question)
		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()

		select {
		case <-tObj.C:
			fmt.Println()
			break problemLoop
		case iAns := <-ansC:
			if iAns == p.answer {
				correctAnswer++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}
	}
	// 7. loop through the problems, print the questions , we'll accept the answers
	// 8. we'll calculate the output and print out the result
	fmt.Printf("Your result is %d out of %d\n", correctAnswer, len(problems))
	fmt.Printf("Press enter to exit")
	<-ansC
}
