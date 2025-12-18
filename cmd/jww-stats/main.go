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
	// ODA FileConverter results
	ODAWarnings int
	ODAErrors   int
	ODAStatus   string
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

	// Print ODA FileConverter results table
	fmt.Println()
	fmt.Println("## ODA FileConverter Results")
	fmt.Println()
	fmt.Println("| File | Warnings | Errors | Status |")
	fmt.Println("|------|----------|--------|--------|")

	for _, s := range allStats {
		fmt.Printf("| `%s` | %d | %d | %s |\n",
			filepath.Base(s.Name), s.ODAWarnings, s.ODAErrors, s.ODAStatus)
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
	totalEzdxfFixes := 0
	odaPassFiles := 0
	for _, s := range allStats {
		if s.Error == "" {
			successFiles++
			if s.DXFError == "" {
				dxfSuccessFiles++
				totalEzdxfFixes += s.EzdxfFixes
				if s.EzdxfErrors == 0 {
					ezdxfPassFiles++
				}
				if s.ODAErrors == 0 {
					odaPassFiles++
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
	fmt.Printf("- ezdxf total fixes applied: %d\n", totalEzdxfFixes)
	fmt.Printf("- ODA FileConverter passed (0 errors): %d\n", odaPassFiles)
}

func parseFile(path string) FileStats {
	stats := FileStats{Name: path, EzdxfStatus: "⏭️ Skipped", ODAStatus: "⏭️ Skipped"}

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

	// Run ODA FileConverter
	odaWarnings, odaErrors, odaStatus := runODAFileConverter(tmpPath)
	stats.ODAWarnings = odaWarnings
	stats.ODAErrors = odaErrors
	stats.ODAStatus = odaStatus

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

// runODAFileConverter runs ODA FileConverter on a DXF file and parses the results.
func runODAFileConverter(dxfPath string) (warnings, errors int, status string) {
	// Create temporary directories for input and output
	tmpDir, err := os.MkdirTemp("", "oda-input-*")
	if err != nil {
		return 0, 0, "⏭️ temp dir error"
	}
	defer os.RemoveAll(tmpDir)

	outDir, err := os.MkdirTemp("", "oda-output-*")
	if err != nil {
		return 0, 0, "⏭️ temp dir error"
	}
	defer os.RemoveAll(outDir)

	// Copy DXF file to input directory
	dxfContent, err := os.ReadFile(dxfPath)
	if err != nil {
		return 0, 0, "⏭️ read error"
	}
	inputPath := filepath.Join(tmpDir, "input.dxf")
	if err := os.WriteFile(inputPath, dxfContent, 0644); err != nil {
		return 0, 0, "⏭️ write error"
	}

	// Run ODAFileConverter
	// Arguments: <input_dir> <output_dir> <output_version> <output_format> <recursive> <audit>
	cmd := exec.Command("ODAFileConverter", tmpDir, outDir, "ACAD2018", "DWG", "0", "1")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		// Check if it's a "command not found" type error
		if strings.Contains(err.Error(), "executable file not found") {
			return 0, 0, "⏭️ ODA not available"
		}
	}

	// Look for .err file in output directory
	errFiles, _ := filepath.Glob(filepath.Join(outDir, "*.err"))
	if len(errFiles) == 0 {
		// Check if DWG was created successfully
		dwgFiles, _ := filepath.Glob(filepath.Join(outDir, "*.dwg"))
		if len(dwgFiles) > 0 {
			return 0, 0, "✅"
		}
		return 0, 1, "❌ no output"
	}

	// Parse error file
	errContent, _ := os.ReadFile(errFiles[0])
	lines := strings.Split(string(errContent), "\n")

	for _, line := range lines {
		if strings.Contains(line, "ODA Warning:") {
			warnings++
		}
		if strings.Contains(line, "OdError") || strings.Contains(line, "ODA Error:") {
			errors++
		}
	}

	if errors > 0 {
		return warnings, errors, fmt.Sprintf("❌ %d errors", errors)
	}
	if warnings > 0 {
		return warnings, errors, fmt.Sprintf("⚠️ %d warnings", warnings)
	}
	return 0, 0, "✅"
}
