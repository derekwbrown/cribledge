package filereader

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	//"math"
)

type filestate struct {
	filename string // the name of the file
	fh	*os.File // the file handle
	lastEndPos int64 // the last end position
	lastStartPos int64 // the last start position
	readsize int64
	overflow []byte
	readbuf []byte
}
var ( 
	readIncrement = int64(1024)
)
// MatchFunc will be called for every matching line in the file.
// if the function returns false, reading will stop independent of 
// whether the count has been reached
type MatchFunc func(string) bool

// ReadFile reads the file from the end, and (optionally) matches each line
func ReverseReadFile(filename string, matchcount int, matchstring string, mf MatchFunc) error {
	// open the file
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	// get the file size
	fi, err := fh.Stat()
	if err != nil {
		return err
	}
	fs := fi.Size()

	state := &filestate{
		filename: filename,
		fh: fh,
		lastEndPos: fs,
		lastStartPos: fs,
		readsize: readIncrement,
		readbuf: make([]byte, readIncrement),
	}

	matches := 0

	for {
		readblocksize := min(readIncrement, state.lastEndPos)
		var offset int64
		if state.lastStartPos > readIncrement {
			offset = state.lastStartPos - readIncrement
			state.lastStartPos = offset
		}
		n, err := fh.ReadAt(state.readbuf[:readblocksize], offset)
		if err != nil {
			return err
		}
		if n == 0 {
			return nil
		}
		state.lastEndPos -= int64(n)

		// now split the buffer into lines
		lines := bytes.Split(state.readbuf[:n], []byte{'\n'})

		if state.overflow != nil {
			lines[len(lines)-1] = append(lines[len(lines)-1], state.overflow...)
			state.overflow = nil
		}
		// stop at 1, because we need to check on the next read if
		// the first line was complete
		for i := len(lines) - 1; i >= 1; i-- {
			if len(lines[i]) == 0 {
				continue
			}
			sent, more := sendOnMatch(string(lines[i]), matchstring, mf)
			if sent {
				matches++
			}
			
			if !more {
				return fmt.Errorf("callback aborted")
			}
			// got all the lines we wanted
			if matchcount > 0 && matches >= matchcount {
				return nil
			}
		}
		if state.lastEndPos == 0{
			sendOnMatch(string(lines[0]), matchstring, mf)
			break
		}
		state.overflow = make([]byte, len(lines[0]))
		copy(state.overflow, lines[0])
	}
	return nil
}

func sendOnMatch(line string, matchstring string, mf MatchFunc) (sent, stop bool) {
	if matchstring == "" || strings.Contains(string(line), matchstring) {
		return true, mf(strings.TrimRight(string(line), "\r\n")) 
	}
	return false, true
}
