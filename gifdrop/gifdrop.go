package main

import (
	"flag"
	"fmt"
	"image/gif"
	"log"
	"os"
)

var skip, threshold int

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [<option> ...] <infile> <outfile>\n",
		os.Args[0])
	fmt.Fprintln(os.Stderr, `
Reduce the size of an animated GIF by dropping frames.
`)
	fmt.Fprintln(os.Stderr, "Options:")
	flag.PrintDefaults()
}

func parseFlags() {
	flag.Usage = usage
	flag.IntVar(&skip, "skip", 1,
		"number of frames to merge into a single frame")
	flag.IntVar(&threshold, "threshold", 100,
		"frames with delays >= this number are always kept")
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
}

func readGIF(filename string) (*gif.GIF, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return gif.DecodeAll(file)
}

func writeGIF(g *gif.GIF, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gif.EncodeAll(file, g); err != nil {
		return err
	}

	return nil
}

func reduceGIF(g *gif.GIF) {
	skipped, total := 0, len(g.Image)
	for i := 1; i < len(g.Image)-1; i++ {
		for j := 0; j < skip-1 && i < len(g.Image)-1; j++ {
			if g.Delay[i] >= threshold {
				break
			}
			g.Image = append(g.Image[:i], g.Image[i+1:]...)
			if g.Delay[i] >= 2 {
				g.Delay[i-1] += g.Delay[i]
			} else {
				g.Delay[i-1] += 10
			}
			g.Delay = append(g.Delay[:i], g.Delay[i+1:]...)
			skipped++
		}
	}
	fmt.Printf("Skipped %d of %d frames (%.0f%%)\n", skipped, total,
		100*float64(skipped)/float64(total))
}

func main() {
	log.SetFlags(log.Lshortfile)

	parseFlags()

	g, err := readGIF(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	reduceGIF(g)

	if err := writeGIF(g, flag.Arg(1)); err != nil {
		log.Fatal(err)
	}
}
