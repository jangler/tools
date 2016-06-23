package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// command-line flags
var (
	echo  = flag.Bool("echo", true, "print response to standard output")
	field = flag.String("field", "file", "name of form field")
)

// die prints args to stderr and exits with nonzero status
func die(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

func main() {
	// handle command-line args
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [url] [file]\n\n", os.Args[0])
		fmt.Fprint(os.Stderr, "Send a file over a HTTP POST request.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}
	url := flag.Arg(0)
	path := flag.Arg(1)

	// make the request, in the sense of constructing it
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	filename := filepath.Base(path)
	fileWriter, err := writer.CreateFormFile(*field, filename)
	if err != nil {
		die(err)
	}
	f, err := os.Open(path)
	if err != nil {
		die(err)
	}
	io.Copy(fileWriter, f)
	f.Close()
	if err := writer.Close(); err != nil {
		die(err)
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		die(err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// make the request, in the sense of using it
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		die(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		die(resp.Status)
	}

	// echo response if desired
	if *echo {
		io.Copy(os.Stdout, resp.Body)
		fmt.Println()
	}
}
