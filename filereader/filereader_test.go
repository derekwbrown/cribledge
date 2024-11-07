package filereader

import (
	"fmt"
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

var expectedLengths = []int{ 38, 36, 34, 32 }

func TestReverseReadFileCanReadWhole(t *testing.T) {
	
	testinfile := filepath.Join("testdata", "simpleinput.txt")
	linecount := 0
	err := ReverseReadFile(testinfile, 0, "", func(line string) bool {
		assert.Equal(t, expectedLengths[linecount], len(line))
		linecount++
		fmt.Println(line)
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
		fmt.Println(line)
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
				fmt.Println(line)
				return true
			})
			assert.Nil(t, err)
			assert.Equal(t, 4, linecount)
		})
	}
}