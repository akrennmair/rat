package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"syscall"
)

func extractArchive() int {
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
			fmt.Fprintf(os.Stderr, "%s\n", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeReg, tar.TypeRegA:
			var f *os.File
			if f, err = os.OpenFile(hdr.Name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.FileMode(hdr.Mode)); err == nil {
				io.Copy(f, archive)
				f.Close()
				err = os.Chtimes(hdr.Name, hdr.ModTime, hdr.ModTime)
			}
		case tar.TypeDir:
			err = os.Mkdir(hdr.Name, os.FileMode(hdr.Mode))
			if err != nil {
				patherr, ok := err.(*os.PathError)
				if ok && patherr.Err == syscall.EEXIST {
					err = nil
				}
			}
		case tar.TypeSymlink:
			// TODO: implement!
		case tar.TypeFifo:
			err = syscall.Mkfifo(hdr.Name, uint32(hdr.Mode))
		case tar.TypeChar:
			fmt.Fprintf(os.Stderr, "Sorry, character devices not supported yet\n")
			//err = syscall.Mknod(hdr.Name, os.S_IFCHR, 
		case tar.TypeBlock:
			fmt.Fprintf(os.Stderr, "Sorry, block devices not supported yet\n")
			//err = syscall.Mknod(hdr.Name, os.S_IFBLK,
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	return 0
}
