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

var itsFlag bool

func parseFlags() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <file>...\n", os.Args[0])
		fmt.Fprintln(os.Stderr, description)
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.BoolVar(&itsFlag, "its", itsFlag, "dump in ITS format instead of WAV")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}
}

func dumpSample(source string, index int, sample *impulse.Sample) error {
	ext := "wav"
	if itsFlag {
		ext = "its"
	}
	filename := fmt.Sprintf("%s-%03d.%s", source[:len(source)-3], index, ext)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if itsFlag {
		return sample.Write(file)
	} else {
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
		return waveFile.Write(file)
	}
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
