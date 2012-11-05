package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
)

func listArchive() {
	archive := tar.NewReader(input)

	for {
		hdr, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured while reading archive: %v\n", err)
			os.Exit(1)
		}

		if verbose {
			mode := os.FileMode(hdr.Mode)
			switch hdr.Typeflag {
			case tar.TypeDir:
				mode |= os.ModeDir
			case tar.TypeSymlink:
				mode |= os.ModeSymlink
			case tar.TypeFifo:
				mode |= os.ModeNamedPipe
			case tar.TypeChar:
				mode |= os.ModeDevice | os.ModeCharDevice
			case tar.TypeBlock:
				mode |= os.ModeDevice
			}
			// TODO: handle uid/gid and/or uname/gname
			fmt.Fprintf(os.Stdout, "%s %9d %s ", mode.String(), hdr.Size, hdr.ModTime.Format("2006-01-02 15:04"))
		}
		fmt.Fprintf(os.Stdout, "%s\n", hdr.Name)
	}

	os.Exit(0)
}
