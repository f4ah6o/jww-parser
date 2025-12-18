// Command jww-stats collects entity statistics from JWW files.
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/f4ah6o/jww-dxf/dxf"
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
	// DXF conversion results
	DXFEntities int
	DXFLayers   int
	DXFBlocks   int
	DXFError    string
	// ezdxf audit results
	EzdxfErrors int
	EzdxfFixes  int
	EzdxfStatus string
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

	// Print DXF conversion results table
	fmt.Println()
	fmt.Println("## DXF Conversion Results")
	fmt.Println()
	fmt.Println("| File | DXF Entities | DXF Layers | DXF Blocks | Conversion Status |")
	fmt.Println("|------|--------------|------------|------------|-------------------|")

	for _, s := range allStats {
		status := "✅"
		if s.DXFError != "" {
			status = "❌ " + s.DXFError
		} else if s.Error != "" {
			status = "⏭️ Parse failed"
		}
		fmt.Printf("| `%s` | %d | %d | %d | %s |\n",
			filepath.Base(s.Name), s.DXFEntities, s.DXFLayers, s.DXFBlocks, status)
	}

	// Print ezdxf audit results table
	fmt.Println()
	fmt.Println("## ezdxf Audit Results")
	fmt.Println()
	fmt.Println("| File | Errors | Fixes | Status |")
	fmt.Println("|------|--------|-------|--------|")

	for _, s := range allStats {
		fmt.Printf("| `%s` | %d | %d | %s |\n",
			filepath.Base(s.Name), s.EzdxfErrors, s.EzdxfFixes, s.EzdxfStatus)
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
	dxfSuccessFiles := 0
	ezdxfPassFiles := 0
	for _, s := range allStats {
		if s.Error == "" {
			successFiles++
			if s.DXFError == "" {
				dxfSuccessFiles++
				if s.EzdxfErrors == 0 {
					ezdxfPassFiles++
				}
			}
		} else {
			errorFiles++
		}
	}
	fmt.Printf("- Total files: %d\n", totalFiles)
	fmt.Printf("- Successfully parsed: %d\n", successFiles)
	fmt.Printf("- Parse errors: %d\n", errorFiles)
	fmt.Printf("- Successfully converted to DXF: %d\n", dxfSuccessFiles)
	fmt.Printf("- ezdxf audit passed (0 errors): %d\n", ezdxfPassFiles)
}

func parseFile(path string) FileStats {
	stats := FileStats{Name: path, EzdxfStatus: "⏭️ Skipped"}

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

	// Convert to DXF and collect statistics
	dxfDoc := dxf.ConvertDocument(doc)
	stats.DXFEntities = len(dxfDoc.Entities)
	stats.DXFLayers = len(dxfDoc.Layers)
	stats.DXFBlocks = len(dxfDoc.Blocks)

	// Write DXF to temp file and run ezdxf audit
	tmpFile, err := os.CreateTemp("", "jww-stats-*.dxf")
	if err != nil {
		stats.EzdxfStatus = "❌ temp file error"
		return stats
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	dxfStr := dxf.ToString(dxfDoc)
	if _, err := tmpFile.WriteString(dxfStr); err != nil {
		tmpFile.Close()
		stats.EzdxfStatus = "❌ write error"
		return stats
	}
	tmpFile.Close()

	// Run ezdxf audit
	errors, fixes, status := runEzdxfAudit(tmpPath)
	stats.EzdxfErrors = errors
	stats.EzdxfFixes = fixes
	stats.EzdxfStatus = status

	return stats
}

// runEzdxfAudit runs ezdxf audit on a DXF file and parses the results.
func runEzdxfAudit(dxfPath string) (errors, fixes int, status string) {
	cmd := exec.Command("uvx", "--from", "git+https://github.com/mozman/ezdxf", "ezdxf", "audit", dxfPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	if err != nil {
		// Check if it's a "command not found" type error
		if strings.Contains(err.Error(), "executable file not found") {
			return 0, 0, "⏭️ ezdxf not available"
		}
	}

	// Parse output for errors and fixes
	// Example: "Found 0 errors, applied 3 fixes" or "No errors found."
	errorsRe := regexp.MustCompile(`Found (\d+) errors`)
	fixesRe := regexp.MustCompile(`applied (\d+) fixes`)
	noErrorsRe := regexp.MustCompile(`No errors found`)

	if noErrorsRe.MatchString(output) {
		return 0, 0, "✅"
	}

	if m := errorsRe.FindStringSubmatch(output); len(m) > 1 {
		fmt.Sscanf(m[1], "%d", &errors)
	}

	if m := fixesRe.FindStringSubmatch(output); len(m) > 1 {
		fmt.Sscanf(m[1], "%d", &fixes)
	}

	if errors == 0 {
		return errors, fixes, "✅"
	}

	return errors, fixes, fmt.Sprintf("⚠️ %d errors", errors)
}
