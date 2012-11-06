package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
)

func listArchive() int {
	archive := tar.NewReader(input)

	for {
		hdr, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured while reading archive: %v\n", err)
			return 1
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

			username := hdr.Uname
			if username == "" {
				username = fmt.Sprintf("%d", hdr.Uid)
			}
			group := hdr.Gname
			if group == "" {
				group = fmt.Sprintf("%d", hdr.Gid)
			}

			user_group := username + "/" + group

			fmt.Fprintf(os.Stdout, "%s %-10s %9d %s ", mode.String(), user_group, hdr.Size, hdr.ModTime.Format("2006-01-02 15:04"))
		}
		fmt.Fprintf(os.Stdout, "%s\n", hdr.Name)
	}

	return 0
}
