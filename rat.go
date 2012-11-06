package main

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

type operation int

const (
	INVALID operation = iota
	CREATE  operation = iota
	LIST    operation = iota
	EXTRACT operation = iota
)

var (
	input     io.Reader      = os.Stdin
	output    io.Writer      = os.Stdout
	fileList                 = []string{}
	filename                 = ""
	directory                = "."
	useGzip                  = false
	useBzip2                 = false
	verbose                  = false
	op        operation      = INVALID
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	parseFlags()

	switch op {
	case CREATE:
		if len(fileList) == 0 {
			printFatal("Refusing to create an empty archive")
		}
		if filename != "" {
			if filename == "-" {
				output = os.Stdout
			} else {
				f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Opening archive %s failed: %v\n", filename, err)
					os.Exit(1)
				}
				defer f.Close()
				output = f
			}
		}
		if useGzip {
			output = gzip.NewWriter(output)
		}
		if useBzip2 {
			printFatal("Sorry, bzip2 compression not yet supported.\n")
		}
		err := os.Chdir(directory)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Changing to %s failed: %v\n", directory, err)
			os.Exit(1)
		}
		os.Exit(createArchive())
	case LIST, EXTRACT:
		if filename != "" {
			if filename == "-" {
				input = os.Stdin
			} else {
				f, err := os.Open(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Opening archive %s failed: %v\n", filename, err)
					os.Exit(1)
				}
				defer f.Close()
				input = f
			}
		}
		if useGzip {
			gzipreader, err := gzip.NewReader(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Reading from gzip failed: %v\n", err)
				os.Exit(1)
			}
			input = gzipreader
		}
		if useBzip2 {
			input = bzip2.NewReader(input)
		}
		if op == LIST {
			os.Exit(listArchive())
		} else {
			err := os.Chdir(directory)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Changing to %s failed: %v\n", directory, err)
				os.Exit(1)
			}
			os.Exit(extractArchive())
		}
	default:
		fmt.Printf("Error: no/invalid operation\n")
	}
}

func parseFlags() {
	i := 1
	j := 2

PARSE_LOOP:
	for {

		for idx, c := range os.Args[i] {
			if idx == 0 {
				if c == '-' {
					continue
				} else if i != 1 {
					break PARSE_LOOP
				}
			}

			switch c {
			case 'c':
				op = CREATE
			case 't':
				op = LIST
			case 'x':
				op = EXTRACT
			case 'h':
				usage()
			case 'f':
				if j >= len(os.Args) {
					printFatal("Option -f requires an argument")
				}
				filename = os.Args[j]
				j++
			case 'z':
				if useBzip2 {
					printFatal("Conflicting compression options")
				}
				useGzip = true
			case 'j', 'y':
				if useGzip {
					printFatal("Conflicting compression options")
				}
				useBzip2 = true
			case 'C':
				if j >= len(os.Args) {
					printFatal("Option -C requires an argument")
				}
				directory = os.Args[j]
				j++
			case 'v':
				verbose = true
			}
		}
		i = j
		j++

		if j >= len(os.Args) {
			break
		}
	}
	if i < len(os.Args) {
		fileList = os.Args[i:]
	}
}

func usage() {
	fmt.Fprintf(os.Stderr,
		`usage:  %s [flags <args>] [files | directories]
        %s {-c} [options] [files | directories]
        %s {-t | -x} [options]

	Options:
		-f file
			Read the archive from or write the archive to the specified file.
			The filename can be - for standard input or standard output.
		-C directory
			In c mode, this changes the directory before adding the files to
			the archive. In x mode, this changes the directory after opening
			the archive but before extracting files from the archive.
		-v
			Produce verbose output.
		-z
			Use gzip compression.
		-j
			Use bzip2 compression (only available for decompression).
		-h
			Show this help.

	Examples:
		rat -cf archive.tar foo bar # create archive.tar from files foo and bar.
		rat -tvf archive.tar        # list all files in archive.tar verbosly.
		rat -xf archive.tar         # extract all files from archive.tar.
`, os.Args[0], os.Args[0], os.Args[0])
	os.Exit(1)

}

func printFatal(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(1)
}
