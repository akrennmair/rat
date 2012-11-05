package main

import (
	"path/filepath"
	"archive/tar"
	"os"
	"fmt"
	"io"
)

func createArchive() {
	archive := tar.NewWriter(output)
	exit_value := 0

	for _, f := range fileList {
		err := filepath.Walk(f, func(path string, info os.FileInfo, err error) error {
			var hdr tar.Header
			hdr.Name = path
			hdr.Size = info.Size()
			hdr.Mode = int64(info.Mode())
			hdr.ModTime = info.ModTime()

			if err := archive.WriteHeader(&hdr); err != nil {
				fmt.Fprintf(os.Stderr, "Writing archive header failed: %v\n", err)
				exit_value = 1
			} else {
				if f, err := os.Open(path); err != nil {
					fmt.Fprintf(os.Stderr, "Opening file failed: %v\n", err)
					exit_value = 1
				} else {
					io.Copy(archive, f)
					f.Close()
					archive.Flush()
				}
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured: %v\n", err)
			exit_value = 1
		}
	}
	archive.Close()
	os.Exit(exit_value)
}
