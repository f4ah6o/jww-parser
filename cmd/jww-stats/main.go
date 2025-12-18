// Command jww-stats collects entity statistics from JWW files.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/f4ah6o/jww-dxf/jww"
)

type FileStats struct {
	Name      string
	Version   uint32
	Lines     int
	Arcs      int
	Points    int
	Texts     int
	Solids    int
	Blocks    int
	BlockDefs int
	Unknown   []string
	Error     string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <dir>\n", os.Args[0])
		os.Exit(1)
	}

	dir := os.Args[1]
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && (filepath.Ext(path) == ".jww" || filepath.Ext(path) == ".JWW") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
		os.Exit(1)
	}

	sort.Strings(files)

	var allStats []FileStats

	for _, file := range files {
		stats := parseFile(file)
		allStats = append(allStats, stats)
	}

	// Print markdown table
	fmt.Println("## Test Data Matrix")
	fmt.Println()
	fmt.Println("| File | Version | Line | Arc | Point | Text | Solid | Block | BlockDef | Error |")
	fmt.Println("|------|---------|------|-----|-------|------|-------|-------|----------|-------|")

	for _, s := range allStats {
		errStr := ""
		if s.Error != "" {
			errStr = "❌ " + s.Error
		}
		fmt.Printf("| `%s` | %d | %d | %d | %d | %d | %d | %d | %d | %s |\n",
			filepath.Base(s.Name), s.Version, s.Lines, s.Arcs, s.Points, s.Texts, s.Solids, s.Blocks, s.BlockDefs, errStr)
	}

	// Print unknown entities summary
	unknownMap := make(map[string]int)
	for _, s := range allStats {
		for _, u := range s.Unknown {
			unknownMap[u]++
		}
	}

	if len(unknownMap) > 0 {
		fmt.Println()
		fmt.Println("## Unknown/Unclassified Entities")
		fmt.Println()
		fmt.Println("| Entity Type | Occurrences |")
		fmt.Println("|-------------|-------------|")
		for k, v := range unknownMap {
			fmt.Printf("| `%s` | %d |\n", k, v)
		}
	}

	// Summary
	fmt.Println()
	fmt.Println("## Summary")
	fmt.Println()
	totalFiles := len(allStats)
	successFiles := 0
	errorFiles := 0
	for _, s := range allStats {
		if s.Error == "" {
			successFiles++
		} else {
			errorFiles++
		}
	}
	fmt.Printf("- Total files: %d\n", totalFiles)
	fmt.Printf("- Successfully parsed: %d\n", successFiles)
	fmt.Printf("- Parse errors: %d\n", errorFiles)
}

func parseFile(path string) FileStats {
	stats := FileStats{Name: path}

	f, err := os.Open(path)
	if err != nil {
		stats.Error = err.Error()
		return stats
	}
	defer f.Close()

	doc, err := jww.Parse(f)
	if err != nil {
		stats.Error = err.Error()
		return stats
	}

	stats.Version = doc.Version
	stats.BlockDefs = len(doc.BlockDefs)

	for _, e := range doc.Entities {
		switch e.Type() {
		case "LINE":
			stats.Lines++
		case "ARC", "CIRCLE":
			stats.Arcs++
		case "POINT":
			stats.Points++
		case "TEXT":
			stats.Texts++
		case "SOLID":
			stats.Solids++
		case "BLOCK":
			stats.Blocks++
		default:
			stats.Unknown = append(stats.Unknown, e.Type())
		}
	}

	return stats
}
