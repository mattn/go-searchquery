package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-searchquery"
)

var (
	lineNumber = flag.Bool("n", false, "print line number with output lines")
	invert     = flag.Bool("v", false, "invert match: select non-matching lines")
	count      = flag.Bool("c", false, "only print a count of matching lines")
	quiet      = flag.Bool("q", false, "quiet mode: suppress normal output")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] QUERY [FILE...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nSearch for lines matching a NIP-50 search query.\n")
		fmt.Fprintf(os.Stderr, "\nQuery Syntax:\n")
		fmt.Fprintf(os.Stderr, "  - Multiple terms are treated as implicit AND\n")
		fmt.Fprintf(os.Stderr, "  - Use quotes for phrases: \"hello world\"\n")
		fmt.Fprintf(os.Stderr, "  - Case-insensitive substring matching\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s \"hello world\" file.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s 'cat dog' file1.txt file2.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  cat file.txt | %s 'search term'\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	query := flag.Arg(0)
	files := flag.Args()[1:]

	exitCode := 0

	if len(files) == 0 {
		// Read from stdin
		if !processReader(os.Stdin, query, "") {
			exitCode = 1
		}
	} else {
		// Process each file
		for _, filename := range files {
			f, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", filename, err)
				exitCode = 2
				continue
			}

			prefix := ""
			if len(files) > 1 {
				prefix = filename + ":"
			}

			if !processReader(f, query, prefix) {
				exitCode = 1
			}
			f.Close()
		}
	}

	os.Exit(exitCode)
}

func processReader(r io.Reader, query, prefix string) bool {
	scanner := bufio.NewScanner(r)
	lineNum := 0
	matchCount := 0
	foundMatch := false

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		match, err := searchquery.Match(line, query)
		if err != nil {
			fmt.Fprintf(os.Stderr, "query error: %v\n", err)
			return false
		}

		if *invert {
			match = !match
		}

		if match {
			matchCount++
			foundMatch = true

			if *quiet {
				// In quiet mode, we just need to know if there's a match
				return true
			}

			if *count {
				// Don't print lines in count mode
				continue
			}

			// Print the matching line
			if *lineNumber {
				fmt.Printf("%s%d:%s\n", prefix, lineNum, line)
			} else {
				fmt.Printf("%s%s\n", prefix, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		return false
	}

	if *count {
		fmt.Printf("%s%d\n", prefix, matchCount)
	}

	return foundMatch
}
