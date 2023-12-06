package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"kal/error"
	"kal/scanner"
	"os"
	"strings"
)

func Run(source string) {
	sc := scanner.NewScanner(source)
	tokens := sc.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func RunFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	Run(string(bytes))

	if error.HadError {
		os.Exit(65)
	}
}

func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		Run(strings.TrimSpace(line))
		if error.HadError {
			break
		}
		error.HadError = false
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: kale [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		RunFile(os.Args[1])
	} else {
		RunPrompt()
	}
}
