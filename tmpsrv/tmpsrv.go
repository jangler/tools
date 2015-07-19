package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	once  = flag.Bool("once", false, "exit after serving one file")
	port  = flag.Uint("port", 8080, "host HTTP server on this port")
	quiet = flag.Bool("quiet", false, "print nothing to standard output")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [<option>]... <file>...\n\n",
			os.Args[0])
		fmt.Fprint(os.Stderr,
			"Serve files specified on the command line over HTTP.\n\n")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, arg := range flag.Args() {
		http.Handle("/"+arg, serverFunc(arg))
	}

	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func serverFunc(path string) http.HandlerFunc {
	f := func(w http.ResponseWriter, req *http.Request) {
		if !*quiet {
			fmt.Println(req.RemoteAddr, req.Method, req.RequestURI)
		}
		http.ServeFile(w, req, path)
		if *once {
			os.Exit(0)
		}
	}
	return http.HandlerFunc(f)
}
