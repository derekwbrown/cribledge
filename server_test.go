package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRequest(filename string, lines int, matchparam string) (*http.Request, error) {
	req, err := http.NewRequest("GET", "http://localhost:8080/getlog", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if filename != "" {
		q.Add("filename", filename)
	}
	if lines > 0 {
		q.Add("count", fmt.Sprintf("%d", lines))
	}
	if matchparam != "" {
		q.Add("match", matchparam)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}
func validateFileReversed(t *testing.T, filename string, body []byte) {
	// we should get the entire file, but in reverse order.

	file, err := os.ReadFile(filepath.Join(rootDirectory, filename))
	require.NoError(t, err)
	filelines := strings.Split(string(file), "\n")
	bodylines := strings.Split(string(body), "\n")

	sourceindex := 0
	bodyindex := len(bodylines) - 1
	for sourceindex < len(filelines) && bodyindex >= 0 {
		if len(filelines[sourceindex]) == 0 {
			sourceindex++
			continue
		}
		if len(bodylines[bodyindex]) == 0 {
			bodyindex--
			continue
		}
		assert.Equal(t, strings.TrimRight(filelines[sourceindex], "\r\n"), strings.TrimRight(bodylines[bodyindex], "\r\n"))
		sourceindex++
		bodyindex--
	}
	assert.Equal(t, -1, bodyindex)
	assert.Equal(t, len(filelines)-1, sourceindex)

}
func TestServer(t *testing.T) {
	s := FileServer()
	defer StopServer(s)
	// override the root directory
	rootDirectory = filepath.Join(".", "testdata")

	t.Run("Test file comes back in reverse order", func(t *testing.T) {
		// ask to see the test `dmesg` file in its entirety
		fn := "dmesg"
		req, err := setupRequest(fn, 0, "")
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		validateFileReversed(t, fn, body)

	})
	t.Run("Test correct number of matches", func(t *testing.T) {
		// ask to see the test `dmesg` file in its entirety
		req, err := setupRequest("dmesg", 0, "cpu")
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		bodylines := strings.Split(string(body), "\n")
		// strip off the last line if it's empty
		if len(bodylines[len(bodylines)-1]) == 0 {
			bodylines = bodylines[:len(bodylines)-1]
		}

		assert.Equal(t, 12, len(bodylines))
	})
	t.Run("Test correct number of matches with match limit", func(t *testing.T) {
		matchlimit := 5
		// ask to see the test `dmesg` file in its entirety
		req, err := setupRequest("dmesg", matchlimit, "cpu")
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		bodylines := strings.Split(string(body), "\n")
		// strip off the last line if it's empty
		if len(bodylines[len(bodylines)-1]) == 0 {
			bodylines = bodylines[:len(bodylines)-1]
		}

		assert.Equal(t, matchlimit, len(bodylines))
	})
	t.Run("Test subdirectory", func(t *testing.T) {
		// ask to see the test `dmesg` file in its entirety
		fn := "path1\\syslog"
		req, err := setupRequest(fn, 0, "")
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		validateFileReversed(t, fn, body)
	})
}
