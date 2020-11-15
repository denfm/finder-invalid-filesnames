package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsValidInDictionaries(t *testing.T) {
	r := [3]bool{true, false, false}
	for k, i := range [3]int32{56, 58, 91} {
		a := IsValidInDictionaries(i)
		assert.True(t, a == r[k])
	}
}

func TestIsWarningInDictionaries(t *testing.T) {
	r := [5]bool{true, true, true, false, false}
	for k, i := range [5]int32{1049, 1053, 1094, 72, 43} {
		a := IsWarningInDictionaries(i)
		assert.True(t, a == r[k])
	}
}

func TestIsValid(t *testing.T) {
	r := [7]bool{true, true, true, true, true, false, false}
	for k, i := range [7]int32{33, 45, 46, 95, 126, 39, 129} {
		a := IsValid(i)
		assert.True(t, a == r[k])
	}
}

func TestValidator(t *testing.T) {
	type s struct {
		word      string
		isValid   bool
		isWarning bool
	}

	data := []s{
		{"привет.txt", true, true},
		{"MCF-tournaments-2011-12.rar", true, false},
		{"heroes-games-2013-boys.pdf", true, false},
		{"stage-moscovia-28-29.2013-regulations.doc", true, false},
		{"ÑÀÎ-Maestro-final-13-14.JPG", false, false},
		{"zdrlet-~!o2013a6.pgn", true, false},
		{"start shkola-2011.xls", false, false},
		{"reg-petrosian-:281212.xls", false, false},
	}

	for _, v := range data {
		isValid, isWarning := Validator(v.word)

		assert.True(t, isValid == v.isValid)
		assert.True(t, isWarning == v.isWarning)
	}
}

func TestScanPath(t *testing.T) {
	path, err := os.Getwd()
	assert.NoError(t, err)

	stats := &stats{
		dirs:  &statsRows{0, 0, 0},
		files: &statsRows{0, 0, 0},
	}
	var buffer []*csvRow

	err = ScanPath(path+"/test_files/", stats, &buffer)
	assert.NoError(t, err)

	assert.Equal(t, stats.dirs.total, 2)
	assert.Equal(t, stats.dirs.warnings, 0)
	assert.Equal(t, stats.dirs.invalids, 0)
	assert.Equal(t, stats.files.total, 3)
	assert.Equal(t, stats.files.warnings, 1)
	assert.Equal(t, stats.files.invalids, 1)
}

func TestSaveCsv(t *testing.T) {
	pathCsv := os.TempDir() + "/finder-testing.csv"
	buffer := []*csvRow{{"FILE", "INVALID", "WARNING", "IS_DIR"}, {"/tmp/hello.txt", "false", "true", "false"}}
	assert.NoError(t, SaveCsv(pathCsv, buffer))

	file, err := os.Open(pathCsv)
	assert.NoError(t, err)

	fi, err := file.Stat()
	assert.NoError(t, err)

	assert.Equal(t, fi.Size(), int64(60))
	assert.NoError(t, os.Remove(pathCsv))
}
