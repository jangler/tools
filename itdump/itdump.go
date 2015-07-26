package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jangler/impulse"
	"github.com/jangler/minipkg/wave"
)

const description = `
Dump all samples from the given IT modules to the working directory.
`

func parseFlags() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <file>...\n", os.Args[0])
		fmt.Fprint(os.Stderr, description)
		os.Exit(2)
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}
}

func dumpSample(source string, index int, sample *impulse.Sample) error {
	waveFile := wave.File{
		Channels:       1,
		SampleRate:     int(sample.Speed),
		BytesPerSample: 1,
		Data:           sample.Data,
	}
	if sample.Flags&impulse.StereoSample != 0 {
		waveFile.Channels = 2
	}
	if sample.Flags&impulse.Quality16Bit != 0 {
		waveFile.BytesPerSample = 2
	}

	filename := fmt.Sprintf("%s-%03d.wav", source[:len(source)-3], index)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return waveFile.Write(file)
}

func dumpFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	module, err := impulse.ReadModule(file)
	if err != nil {
		return err
	}

	for i, sample := range module.Samples {
		if err := dumpSample(filename, i, sample); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %v\n", os.Args[0], filename, err)
		}
	}
	return nil
}

func main() {
	parseFlags()
	for _, arg := range flag.Args() {
		if err := dumpFile(arg); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %v\n", os.Args[0], arg, err)
		}
	}
}
