package main


import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	
)


// query arguments
// filename=name, can include path separators
// count=number, the number of lines to read
// match=string, the string to match

func main() {

	// override the root directory
	rootDirectory = filepath.Join(".", "testdata")

	s := FileServer()
	defer StopServer(s)
	done := make(chan os.Signal, 1)
  	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
  	fmt.Println("Blocking, press ctrl+c to continue...")
  	<-done  // Will block here u
}