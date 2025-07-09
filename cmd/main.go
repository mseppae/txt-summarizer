package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/mseppae/txt-summarizer/summarizer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: summarize <file1> <file2> ...")
		os.Exit(1)
	}

	keySums, err := summarizer.ParseAndSumFiles(os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	keySumPairs := summarizer.SortKeySums(keySums)

	outputFileName := fmt.Sprintf("summary-%d.txt", time.Now().Unix())
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error creating output file %s: %v\n", outputFileName, err)
		os.Exit(1)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	if err := summarizer.WriteSummary(writer, keySumPairs); err != nil {
		fmt.Printf("Error writing to output file %s: %v\n", outputFileName, err)
		os.Exit(1)
	}

	fmt.Printf("Summary written to %s\n", outputFileName)
}
