package bookmarks

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"
)

// Folder represents a bookmark folder containing subfolders or individual bookmarks.
type Folder struct {
	Name     string
	Children []any // Can contain *Folder or Bookmark
}

// Bookmark represents a single bookmarked URL.
type Bookmark struct {
	Title string
	URL   string
}

var (
	folderRe   = regexp.MustCompile(`<H3[^>]*>(.*?)</H3>`)
	bookmarkRe = regexp.MustCompile(`<A[^>]*HREF="([^"]+)"[^>]*>(.*?)</A>`)
)

// ParseHTML reads Netscape HTML bookmark data from an io.Reader and returns a structured root Folder.
func ParseHTML(r io.Reader) (*Folder, error) {
	root := &Folder{Name: "Bookmarks"}
	stack := []*Folder{root}

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// match Folder Creation
		if matches := folderRe.FindStringSubmatch(line); len(matches) > 1 {
			name := html.UnescapeString(matches[1])
			folder := &Folder{Name: name}

			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, folder)

			stack = append(stack, folder)
			continue
		}

		// match Bookmark Link
		if matches := bookmarkRe.FindStringSubmatch(line); len(matches) > 2 {
			url := html.UnescapeString(matches[1])
			title := html.UnescapeString(matches[2])

			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, Bookmark{
				Title: title,
				URL:   url,
			})
			continue
		}

		// Pop stack on folder closure
		if strings.Contains(line, "</DL>") && len(stack) > 1 {
			stack = stack[:len(stack)-1]
		}
	}

	return root, scanner.Err()
}

// WriteMarkdown writes the Folder tree as Markdown to the provided io.Writer.
func WriteMarkdown(w io.Writer, folder *Folder) error {
	return writeMarkdownLevel(w, folder, 0)
}

func writeMarkdownLevel(w io.Writer, folder *Folder, level int) error {
	var err error
	if level == 0 {
		_, err = fmt.Fprintf(w, "# %s\n\n", folder.Name)
	} else {
		heading := strings.Repeat("#", min(level+1, 6))
		_, err = fmt.Fprintf(w, "%s %s\n\n", heading, folder.Name)
	}
	if err != nil {
		return err
	}

	for _, child := range folder.Children {
		switch v := child.(type) {
		case Bookmark:
			if _, err := fmt.Fprintf(w, "- [%s](%s)\n", v.Title, v.URL); err != nil {
				return err
			}

		case *Folder:
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
			if err := writeMarkdownLevel(w, v, level+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
