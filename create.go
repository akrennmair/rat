package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func createArchive() {
	archive := tar.NewWriter(output)
	defer archive.Close()
	exit_value := 0

	for _, f := range fileList {
		err := filepath.Walk(f, func(path string, info os.FileInfo, err error) error {

			var hdr tar.Header
			hdr.Name = path
			hdr.Size = info.Size()
			hdr.Mode = int64(info.Mode())
			hdr.ModTime = info.ModTime()
			// TODO: add uid/gid and/or uname/gname
			// TODO: handle directories, symlinks and other file types

			switch (info.Mode() & os.ModeType) {
			case 0:
				hdr.Typeflag = tar.TypeReg
			case os.ModeDir:
				hdr.Typeflag = tar.TypeDir
			case os.ModeSymlink:
				hdr.Typeflag = tar.TypeSymlink
				linkname, err := os.Readlink(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: can't readlink a symlink\n")
					return nil
				} else {
					hdr.Linkname = linkname
				}
			case os.ModeNamedPipe:
				hdr.Typeflag = tar.TypeFifo
			case os.ModeSocket:
				fmt.Fprintf(os.Stderr, "Warning: can't tar a socket\n")
				return nil
			case os.ModeDevice:
				fmt.Fprintf(os.Stderr, "Warning: device files are currently unsupported\n")
				return nil
			/*
				if (info.Mode() & os.ModeCharDevice) != 0 {
					os.Typeflag = tar.TypeChar
				} else {
					os.Typeflag = tar.TypeBlock
				}
				*/
			}

			if err := archive.WriteHeader(&hdr); err != nil {
				fmt.Fprintf(os.Stderr, "Writing archive header for %s failed: %v\n", path, err)
				exit_value = 1
				return nil
			}
			defer archive.Flush()

			if hdr.Typeflag == tar.TypeReg {
				if f, err := os.Open(path); err != nil {
					fmt.Fprintf(os.Stderr, "Opening file %s failed: %v\n", path, err)
					exit_value = 1
					return nil
				} else {
					io.Copy(archive, f)
					f.Close()
				}
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured: %v\n", err)
			exit_value = 1
		}
	}
	os.Exit(exit_value)
}
