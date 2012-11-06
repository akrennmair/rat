package main

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
)

func makedev(major, minor int64) int {
	return int(((major & 0xfff) << 8) | (minor & 0xff) | ((major &^ 0xfff) << 32) | ((minor & 0xfffff00) << 12))
}

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
			err = os.Link(hdr.Name, hdr.Linkname)
		case tar.TypeFifo:
			err = syscall.Mkfifo(hdr.Name, uint32(hdr.Mode))
		case tar.TypeChar:
			err = errors.New("character devices unsupported")
			err = syscall.Mknod(hdr.Name, syscall.S_IFCHR, makedev(hdr.Devmajor, hdr.Devminor))
		case tar.TypeBlock:
			err = errors.New("block devices unsupported")
			err = syscall.Mknod(hdr.Name, syscall.S_IFBLK, makedev(hdr.Devmajor, hdr.Devminor))
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	return 0
}
