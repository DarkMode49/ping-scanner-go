package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/stoicperlman/fls"
)

type FileIPReader struct{
	file *fls.File
}

// defer file.Close()

func (r *FileIPReader) Open(filePath string) (*fls.File, error) {
	var fileError error
	r.file, fileError = fls.Open(filePath)
	// parent function has the responsibility
	// to close the file
	if fileError != nil {
		return nil, fileError
	}

	return r.file, nil
}

func (r *FileIPReader) ReadIP(line int64) (string, error) {
	// lineBeginningPosition
	_, seekLineError := r.file.SeekLine(line, io.SeekStart)
	if seekLineError != nil {
		return "", seekLineError
	}

	// Seek to the beginning of the line
	// _, seekError := r.file.Seek(lineBeginningPosition, io.SeekStart)
	// if seekError != nil {
	// 	return "", seekError
	// }

	// Use a buffered reader to read the line efficiently
	bufReader := bufio.NewReader(r.file)
	ipLine, ipReadError := bufReader.ReadString('\n')
	if ipReadError != nil && ipReadError != io.EOF {
		return "", ipReadError
	}

	// Remove trailing newline and carriage return if present
	ipLine = strings.TrimRight(ipLine, "\r\n")

	return ipLine, nil
}

// FileIPWriter implements IPWriter to write IPs to a plain text file.
type FileIPWriter struct{}

func (w *FileIPWriter) WriteIPs(filePath string, ips []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, ip := range ips {
		if _, err := fmt.Fprintln(writer, ip); err != nil {
			return fmt.Errorf("failed to write ip %s to file: %w", ip, err)
		}
	}

	return writer.Flush()
}
