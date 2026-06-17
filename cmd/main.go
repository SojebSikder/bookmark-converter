package main

import (
	"fmt"
	"os"

	"github.com/sojebsikder/bookmark-converter/bookmarks"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:")
		fmt.Println("  bc bookmarks.html bookmarks.md")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	if err := convertBookmarks(inputPath, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully wrote Markdown to %s\n", outputPath)
}

func convertBookmarks(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	root, err := bookmarks.ParseHTML(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse bookmarks: %w", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	if err := bookmarks.WriteMarkdown(outputFile, root); err != nil {
		return fmt.Errorf("failed to write markdown: %w", err)
	}

	return nil
}
