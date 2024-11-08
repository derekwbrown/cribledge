package main

import (
	"context"
	"github.com/derekwbrown/cribledge/filereader"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

var (
	rootDirectory = filepath.Join("c:\\", "ProgramData")
)

func FileServer() *http.Server {
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	go func() {
		http.HandleFunc("/getlog", func(w http.ResponseWriter, r *http.Request) {

			filename := r.URL.Query().Get("filename")
			if filename == "" {
				http.Error(w, "filename is required", http.StatusBadRequest)
				return
			}
			var intcount int
			var err error
			count := r.URL.Query().Get("count")
			if count != "" {
				intcount, err = strconv.Atoi(count)
				if err != nil {
					http.Error(w, "count must be a number", http.StatusBadRequest)
					return
				}
			}
			match := r.URL.Query().Get("match")
			completePath := filepath.Join(rootDirectory, filename)

			err = filereader.ReverseReadFile(completePath, intcount, match, func(line string) bool {
				w.Write([]byte(line))
				w.Write([]byte("\r\n"))
				return true
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
		// start the server
		server.ListenAndServe()
	}()
	return server
}

func StopServer(server *http.Server) error {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return nil
}
