// Command jww-stats collects entity statistics from JWW files.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/f4ah6o/jww-dxf/dxf"
	"github.com/f4ah6o/jww-dxf/jww"
)

// Command line flags
var odaFlag = flag.Bool("oda", false, "Run ODA FileConverter check (disabled by default)")

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
	// ezdxf info results (from ezdxf info -s)
	EzdxfInfoEntities int // Entities in modelspace
	EzdxfInfoLayers   int // LAYER table entries
	EzdxfInfoBlocks   int // BLOCK_RECORD table entries
	EzdxfInfoStatus   string
	// ODA FileConverter results
	ODAWarnings int
	ODAErrors   int
	ODAStatus   string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <dir>\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	dir := flag.Arg(0)
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

	// Process files in parallel
	allStats := make([]FileStats, len(files))
	var wg sync.WaitGroup

	for i, file := range files {
		wg.Add(1)
		go func(idx int, filePath string) {
			defer wg.Done()
			allStats[idx] = parseFile(filePath)
		}(i, file)
	}

	wg.Wait()

	// Build Test Data Matrix rows
	var testDataRows [][]string
	for _, s := range allStats {
		errStr := ""
		if s.Error != "" {
			errStr = "❌ " + s.Error
		}
		testDataRows = append(testDataRows, []string{
			"`" + filepath.Base(s.Name) + "`",
			fmt.Sprintf("%d", s.Version),
			fmt.Sprintf("%d", s.Lines),
			fmt.Sprintf("%d", s.Arcs),
			fmt.Sprintf("%d", s.Points),
			fmt.Sprintf("%d", s.Texts),
			fmt.Sprintf("%d", s.Solids),
			fmt.Sprintf("%d", s.Blocks),
			fmt.Sprintf("%d", s.BlockDefs),
			errStr,
		})
	}

	fmt.Println("## Test Data Matrix")
	fmt.Println()
	printTable([]string{"File", "Version", "Line", "Arc", "Point", "Text", "Solid", "Block", "BlockDef", "Error"}, testDataRows)

	// Build DXF Conversion Results rows
	var dxfRows [][]string
	for _, s := range allStats {
		status := "✅"
		if s.DXFError != "" {
			status = "❌ " + s.DXFError
		} else if s.Error != "" {
			status = "⏭️ Parse failed"
		}
		jwwTotal := s.Lines + s.Arcs + s.Points + s.Texts + s.Solids + s.Blocks
		diff := s.DXFEntities - jwwTotal
		diffStr := fmt.Sprintf("%+d", diff)
		if diff == 0 {
			diffStr = "0 ✅"
		}
		dxfRows = append(dxfRows, []string{
			"`" + filepath.Base(s.Name) + "`",
			fmt.Sprintf("%d", jwwTotal),
			fmt.Sprintf("%d", s.DXFEntities),
			diffStr,
			status,
		})
	}

	fmt.Println()
	fmt.Println("## DXF Conversion Results (Entity Count Comparison)")
	fmt.Println()
	printTable([]string{"File", "JWW Entities", "DXF Entities", "Diff", "Status"}, dxfRows)

	// Build ezdxf Audit Results rows
	var auditRows [][]string
	for _, s := range allStats {
		auditRows = append(auditRows, []string{
			"`" + filepath.Base(s.Name) + "`",
			fmt.Sprintf("%d", s.EzdxfErrors),
			fmt.Sprintf("%d", s.EzdxfFixes),
			s.EzdxfStatus,
		})
	}

	fmt.Println()
	fmt.Println("## ezdxf Audit Results")
	fmt.Println()
	printTable([]string{"File", "Errors", "Fixes", "Status"}, auditRows)

	// Build ezdxf Info Results rows
	var infoRows [][]string
	for _, s := range allStats {
		infoRows = append(infoRows, []string{
			"`" + filepath.Base(s.Name) + "`",
			fmt.Sprintf("%d", s.EzdxfInfoEntities),
			fmt.Sprintf("%d", s.EzdxfInfoLayers),
			fmt.Sprintf("%d", s.EzdxfInfoBlocks),
			s.EzdxfInfoStatus,
		})
	}

	fmt.Println()
	fmt.Println("## ezdxf Info Results (DXF File Statistics)")
	fmt.Println()
	printTable([]string{"File", "Entities", "Layers", "Blocks", "Status"}, infoRows)

	// Print ODA FileConverter results table (only if --oda flag is set)
	if *odaFlag {
		var odaRows [][]string
		for _, s := range allStats {
			odaRows = append(odaRows, []string{
				"`" + filepath.Base(s.Name) + "`",
				fmt.Sprintf("%d", s.ODAWarnings),
				fmt.Sprintf("%d", s.ODAErrors),
				s.ODAStatus,
			})
		}

		fmt.Println()
		fmt.Println("## ODA FileConverter Results")
		fmt.Println()
		printTable([]string{"File", "Warnings", "Errors", "Status"}, odaRows)
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
				if *odaFlag && s.ODAErrors == 0 {
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
	if *odaFlag {
		fmt.Printf("- ODA FileConverter passed (0 errors): %d\n", odaPassFiles)
	}
}

func parseFile(path string) FileStats {
	odaStatus := "⏭️ Disabled"
	if *odaFlag {
		odaStatus = "⏭️ Skipped"
	}
	stats := FileStats{Name: path, EzdxfStatus: "⏭️ Skipped", EzdxfInfoStatus: "⏭️ Skipped", ODAStatus: odaStatus}

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

	// Run ezdxf info
	runEzdxfInfo(tmpPath, &stats)

	// Run ODA FileConverter (only if --oda flag is set)
	if *odaFlag {
		odaWarnings, odaErrors, odaStatus := runODAFileConverter(tmpPath)
		stats.ODAWarnings = odaWarnings
		stats.ODAErrors = odaErrors
		stats.ODAStatus = odaStatus
	}

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

// runEzdxfInfo runs ezdxf info on a DXF file and parses summary statistics.
func runEzdxfInfo(dxfPath string, stats *FileStats) {
	cmd := exec.Command("uvx", "--from", "git+https://github.com/mozman/ezdxf", "ezdxf", "info", "-s", dxfPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			stats.EzdxfInfoStatus = "⏭️ ezdxf not available"
			return
		}
	}

	// Parse summary statistics from ezdxf info -s output
	// Example output format:
	// Entities in modelspace: 695
	// LAYER table entries: 258
	// BLOCK_RECORD table entries: 2

	// Parse entities in modelspace
	entitiesRe := regexp.MustCompile(`Entities in modelspace:\s*(\d+)`)
	if m := entitiesRe.FindStringSubmatch(output); len(m) > 1 {
		fmt.Sscanf(m[1], "%d", &stats.EzdxfInfoEntities)
	}

	// Parse layer table entries
	layersRe := regexp.MustCompile(`LAYER table entries:\s*(\d+)`)
	if m := layersRe.FindStringSubmatch(output); len(m) > 1 {
		fmt.Sscanf(m[1], "%d", &stats.EzdxfInfoLayers)
	}

	// Parse block record table entries
	blocksRe := regexp.MustCompile(`BLOCK_RECORD table entries:\s*(\d+)`)
	if m := blocksRe.FindStringSubmatch(output); len(m) > 1 {
		fmt.Sscanf(m[1], "%d", &stats.EzdxfInfoBlocks)
	}

	stats.EzdxfInfoStatus = "✅"
}

// printTable prints a markdown table with aligned columns.
// headers is a slice of column header strings.
// rows is a slice of row data, where each row is a slice of cell strings.
func printTable(headers []string, rows [][]string) {
	// Calculate column widths (using rune width for Unicode support)
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = runeWidth(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				w := runeWidth(cell)
				if w > widths[i] {
					widths[i] = w
				}
			}
		}
	}

	// Print header
	fmt.Print("|")
	for i, h := range headers {
		fmt.Printf(" %-*s |", widths[i]+runeWidth(h)-len(h), h)
	}
	fmt.Println()

	// Print separator
	fmt.Print("|")
	for _, w := range widths {
		fmt.Print(strings.Repeat("-", w+2) + "|")
	}
	fmt.Println()

	// Print rows
	for _, row := range rows {
		fmt.Print("|")
		for i, cell := range row {
			if i < len(widths) {
				// Pad with spaces accounting for rune width
				padding := widths[i] - runeWidth(cell) + len(cell)
				fmt.Printf(" %-*s |", padding, cell)
			}
		}
		fmt.Println()
	}
}

// runeWidth returns the display width of a string, accounting for wide characters.
func runeWidth(s string) int {
	width := 0
	for _, r := range s {
		// Wide characters (CJK, emoji, etc.) take 2 columns
		if r >= 0x1100 && (r <= 0x115F || // Hangul Jamo
			r == 0x2329 || r == 0x232A || // Angle brackets
			(r >= 0x2E80 && r <= 0xA4CF && r != 0x303F) || // CJK
			(r >= 0xAC00 && r <= 0xD7A3) || // Hangul Syllables
			(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility Ideographs
			(r >= 0xFE10 && r <= 0xFE19) || // Vertical forms
			(r >= 0xFE30 && r <= 0xFE6F) || // CJK Compatibility Forms
			(r >= 0xFF00 && r <= 0xFF60) || // Fullwidth Forms
			(r >= 0xFFE0 && r <= 0xFFE6) || // Fullwidth Forms
			(r >= 0x1F300 && r <= 0x1F9FF) || // Emoji
			(r >= 0x20000 && r <= 0x2FFFF)) { // CJK Extension B+
			width += 2
		} else {
			width += 1
		}
	}
	return width
}
