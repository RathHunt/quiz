package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func timeout(t chan bool, timeLimit int) {
	time.Sleep(time.Duration(timeLimit) * time.Second)
	t <- false
}

func readAns(s chan string) {
	in := bufio.NewReader(os.Stdin)
	answer, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	s <- answer
}

func equals(a, b string) bool {
	return strings.ToLower(strings.Trim(a, "\r\n")) == strings.ToLower(b)
}

func main() {

	var filename = flag.String("filename", "problems.csv", "file containing quiz")
	var timelimit = flag.Int("timelimit", 30, "time limit in seconds")
	var shuffle = flag.Bool("shuffle", false, "shuffle the quiz order")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press Enter to start timer")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	t := make(chan bool)
	go timeout(t, *timelimit)

	records, _ := csv.NewReader(f).ReadAll()

	if *shuffle {
		for i := range records {
			j := rand.Intn(i + 1)
			records[i], records[j] = records[j], records[i]
		}
	}

	var total = len(records)

	var correct int

loop:
	for _, record := range records {

		fmt.Println(record[0])

		input := make(chan string)

		go readAns(input)

		select {
		case <-t:
			fmt.Println("Timeout.")
			break loop
		case answer := <-input:
			if equals(answer, record[1]) {
				correct++
			}

		}
	}

	fmt.Printf("Correct: %d\nTotal: %d", correct, total)

}
