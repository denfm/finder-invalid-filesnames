package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

type statsRows struct {
	total    int
	invalids int
	warnings int
}

type stats struct {
	estimatedTimeLabel string
	dirs               *statsRows
	files              *statsRows
}

type csvRow struct {
	file    string
	invalid string
	warning string
	isDir   string
}

var parsePath, outputCsvPath string

var ValidDictionaries = [3][2]int32{
	{48, 57},  // 0-9
	{65, 90},  // A-Z
	{97, 122}, // a-z
}

var WarningDictionaries = [1][2]int32{
	{1040, 1103}, // А-Я|а-я
}

var ValidCodes = [5]int32{
	33,  // !
	45,  // -
	46,  // .
	95,  // _
	126, // ~
}

// https://www.tamasoft.co.jp/en/general-info/unicode-decimal.html
func main() {
	flag.StringVar(&parsePath, "parse-path", "", "The directory to be scanned recursively")
	flag.StringVar(&outputCsvPath, "output-csv-path", "/tmp/finder-invalid-filenames.csv",
		"The path to the csv file for the result")

	flag.Parse()

	if parsePath == "" {
		log.Fatal("Please specify path for recursive scanning")
	}

	stats := &stats{
		dirs:  &statsRows{0, 0, 0},
		files: &statsRows{0, 0, 0},
	}
	buffer := []*csvRow{{"FILE", "INVALID", "WARNING", "IS_DIR"}}

	err := ScanPath(parsePath, stats, &buffer)

	if err != nil {
		log.Fatal(err)
	}

	err = SaveCsv(outputCsvPath, buffer)

	if err != nil {
		log.Fatal(err)
	}

	headerStatsRow := "Type\tInvalid\tWarnings\tTotal\t"
	filesStatsRow := fmt.Sprintf("Files\t%d\t%d\t%d\t", stats.files.invalids, stats.files.warnings, stats.files.total)
	dirStatsRow := fmt.Sprintf("Dirs\t%d\t%d\t%d\t", stats.dirs.invalids, stats.dirs.warnings, stats.dirs.total)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, '.', tabwriter.AlignRight)
	_, _ = fmt.Fprintln(w, headerStatsRow)
	_, _ = fmt.Fprintln(w, filesStatsRow)
	_, _ = fmt.Fprintln(w, dirStatsRow)
	_ = w.Flush()
}

func SaveCsv(path string, buffer []*csvRow) error {
	if len(buffer) > 1 {
		_ = os.Remove(path)

		file, err := os.Create(path)
		if err != nil {
			return err
		}

		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		for _, value := range buffer {
			err := writer.Write([]string{value.file, value.isDir, value.invalid, value.warning})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ScanPath(path string, stats *stats, buffer *[]*csvRow) error {
	return filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		isValid, isWarning := Validator(fi.Name())

		if fi.IsDir() {
			stats.dirs.total++

			if !isValid {
				stats.dirs.invalids++
			}
			if isWarning {
				stats.dirs.warnings++
			}
		} else {
			stats.files.total++

			if !isValid {
				stats.files.invalids++
			}
			if isWarning {
				stats.files.warnings++
			}
		}

		if !isValid || isWarning {
			*buffer = append(*buffer, &csvRow{
				strings.TrimLeft(strings.TrimPrefix(file, parsePath), "/"),
				strconv.FormatBool(!isValid),
				strconv.FormatBool(isWarning),
				strconv.FormatBool(fi.IsDir()),
			})
		}
		return nil
	})
}

func IsValidInDictionaries(unicode int32) bool {
	for _, vCode := range ValidDictionaries {
		if unicode >= vCode[0] && unicode <= vCode[1] {
			return true
		}
	}

	return false
}

func IsWarningInDictionaries(unicode int32) bool {
	for _, vCode := range WarningDictionaries {
		if unicode >= vCode[0] && unicode <= vCode[1] {
			return true
		}
	}

	return false
}

func IsValid(unicode int32) bool {
	for _, vCode := range ValidCodes {
		if unicode == vCode {
			return true
		}
	}

	return false
}

func Validator(str string) (bool, bool) {
	uniCodes := []rune(str)

	isWarning := false
	isInvalid := false

	for _, unicode := range uniCodes {
		if IsWarningInDictionaries(unicode) {
			isWarning = true
		} else if !IsValidInDictionaries(unicode) && !IsValid(unicode) {
			isInvalid = true
			break
		}
	}

	return !isInvalid, isWarning
}
