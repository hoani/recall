package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("log.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	formatter := NewLineFormatter()
	for scanner.Scan() {
		line, err := formatter.Format(scanner.Bytes())
		if err != nil {
			panic(err)
		}
		fmt.Println(line)
	}
}
