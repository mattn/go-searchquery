package main

import (
	"bufio"
	"bytes"
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

func usage() {
	fmt.Fprint(os.Stderr, "Usage: query [OPTIONS] QUERY [FILE...]\n")
	fmt.Fprint(os.Stderr, "\nSearch for lines matching a NIP-50 search query.\n")
	fmt.Fprint(os.Stderr, "\nQuery Syntax:\n")
	fmt.Fprint(os.Stderr, "  - Multiple terms are treated as implicit AND\n")
	fmt.Fprint(os.Stderr, "  - Use quotes for phrases: \"hello world\"\n")
	fmt.Fprint(os.Stderr, "  - Case-insensitive substring matching\n")
	fmt.Fprint(os.Stderr, "\nExamples:\n")
	fmt.Fprint(os.Stderr, "  query \"hello world\" file.txt\n")
	fmt.Fprint(os.Stderr, "  query 'cat dog' file1.txt file2.txt\n")
	fmt.Fprint(os.Stderr, "  cat file.txt | query 'search term'\n")
	fmt.Fprint(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func isBinary(buf []byte) bool {
	return bytes.IndexByte(buf, 0) != -1
}

func process(query string, files []string) int {
	exitCode := 0
	for _, filename := range files {
		if st, err := os.Stat(filename); err == nil && st.IsDir() {
			ff, err := os.ReadDir(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", filename, err)
				exitCode = 2
				continue
			}
			for _, entry := range ff {
				process(query, []string{filename + "/" + entry.Name()})
			}
		} else {
			f, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", filename, err)
				exitCode = 2
				continue
			}

			buf := make([]byte, 512)
			n, err := f.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "%s: %v\n", filename, err)
				f.Close()
				exitCode = 2
				continue
			}

			if isBinary(buf[:n]) {
				f.Close()
				continue
			}

			prefix := filename + ":"
			if !processReader(io.MultiReader(bytes.NewReader(buf[:n]), f), query, prefix) {
				exitCode = 1
			}
			f.Close()
		}
	}

	return exitCode
}

func main() {
	flag.Usage = usage

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	query := flag.Arg(0)
	files := flag.Args()[1:]

	if len(files) == 0 {
		// Read from stdin
		if !processReader(os.Stdin, query, "") {
			os.Exit(1)
		}
	} else {
		// Process each file
		os.Exit(process(query, files))
	}
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

	if err := scanner.Err(); err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		return false
	}

	if *count {
		fmt.Printf("%s%d\n", prefix, matchCount)
	}

	return foundMatch
}
