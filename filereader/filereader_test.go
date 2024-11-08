package filereader

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	//"github.com/stretchr/testify/require"
)

var expectedLengths = []int{38, 36, 34, 32}

func TestReverseReadFileCanReadWhole(t *testing.T) {

	testinfile := filepath.Join("testdata", "simpleinput.txt")
	linecount := 0
	err := ReverseReadFile(testinfile, 0, "", func(line string) bool {
		assert.Equal(t, expectedLengths[linecount], len(line))
		linecount++
		return true
	})
	assert.Nil(t, err)
	assert.Equal(t, 4, linecount)
}

func TestReverseReadShortReads(t *testing.T) {
	readIncrement = int64(30)
	testinfile := filepath.Join("testdata", "simpleinput.txt")
	linecount := 0
	err := ReverseReadFile(testinfile, 0, "", func(line string) bool {
		assert.Equal(t, expectedLengths[linecount], len(line))
		linecount++
		return true
	})
	assert.Nil(t, err)
	assert.Equal(t, 4, linecount)
}
func TestReverseReadOnBoundaries(t *testing.T) {
	// just use the expected length for each line as the read length
	for _, l := range expectedLengths {
		t.Run(fmt.Sprintf("Read %d bytes", l), func(t *testing.T) {
			readIncrement = int64(l)
			testinfile := filepath.Join("testdata", "simpleinput.txt")
			linecount := 0
			err := ReverseReadFile(testinfile, 0, "", func(line string) bool {
				assert.Equal(t, expectedLengths[linecount], len(line))
				linecount++
				return true
			})
			assert.Nil(t, err)
			assert.Equal(t, 4, linecount)
		})
	}
}

func TestReverseReadStopsAtCount(t *testing.T) {
	knownlinesinfile := 4
	for i := 1; i < 6; i++ {
		t.Run(fmt.Sprintf("Stop at %d lines", i), func(t *testing.T) {
			testinfile := filepath.Join("testdata", "simpleinput.txt")
			linecount := 0
			err := ReverseReadFile(testinfile, i, "", func(line string) bool {
				assert.Equal(t, expectedLengths[linecount], len(line))
				linecount++
				return true
			})
			assert.Nil(t, err)
			// the test will ask for more lines that are in present at the
			// end of the loop, but we should still only get the max number
			// of lines
			assert.Equal(t, min(i, knownlinesinfile), linecount)
		})
	}
}

func TestReverseReadWithSimpleMatching(t *testing.T) {
	// map of string we're to match, with expected match values
	matchtests := map[string]int{
		"aaaaaa": 1,
		"bbbbbb": 1,
		"cccccc": 1,
		"aba":    0,
		"ddda":   0,
		// this should _not_ match because it's longer than the input string
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa":       0,
		"dddddddddddddddddddddddddddddddddddddddd": 0,
	}
	for matchstring, expectedmatches := range matchtests {
		t.Run(fmt.Sprintf("Match %s", matchstring), func(t *testing.T) {
			testinfile := filepath.Join("testdata", "simpleinput.txt")
			matches := 0
			err := ReverseReadFile(testinfile, 0, matchstring, func(line string) bool {
				matches++
				return true
			})
			assert.Nil(t, err)
			assert.Equal(t, expectedmatches, matches)
		})
	}
}
