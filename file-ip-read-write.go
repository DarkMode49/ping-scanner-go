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
	_, seekLineError := r.file.SeekLine(line, io.SeekStart)
	if seekLineError != nil {
		return "", seekLineError
	}

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

const WriteFilePermissions = 0o644 // owner read/write, others read

func NewFileIPWriter() *FileIPWriter {
	return &FileIPWriter{}
}

func (w *FileIPWriter) Open(filePath string) (io.ReadWriteCloser, error) {
	file, openFileError := os.OpenFile(
		filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		WriteFilePermissions,
	)
	// parent function has the responsibility
	// to close the file

	if openFileError != nil {
		return nil, fmt.Errorf(
			"could not open write file %s! more: %w", filePath, openFileError,
		)
	}

	return file, nil
}

func (w *FileIPWriter) GetBuffer(file io.ReadWriteCloser) (WriterFlusher) {
	return bufio.NewWriter(file)
}

func (w *FileIPWriter) WriteIP(writer WriterFlusher, ip string) error {
	_, writeError := fmt.Fprintln(writer, ip);

	if writeError != nil {
		return writeError
	}

	return writer.Flush()
}
