package main

import (
	"fmt"
	"os"

	"github.com/sojebsikder/bookmark-converter/bookmarks"
)

const version = "v0.1.0"

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			fmt.Printf("bookmark-converter version %s\n", version)
			os.Exit(0)
		}
	}

	if len(os.Args) != 3 {
		fmt.Println("Usage:")
		fmt.Println("  bc <input.html> <output.md>  Convert bookmarks")
		fmt.Println("  bc version                   Show version information")
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
