package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	var link = regexp.MustCompile(`^(.*) (-+) (.*)$`)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("graph G {")
	for scanner.Scan() {
		elements := link.FindStringSubmatch(scanner.Text())
		if len(elements) != 4 {
			log.Fatal("bad entry", scanner.Text())
		}
		fmt.Printf(`"%v" -- "%v" [minlen=%v]`, elements[1], elements[3], len(elements[2]))
	}
	fmt.Println("}")
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
