package main

import (
	"bytes"
	"io"
	"os"
)

func divisionBoundaries(length int, parts int) [][]int {
	if parts <= 0 || length <= 0 {
		return [][]int{}
	}

	offsets := make([][]int, parts)
	base := length / parts
	start := 0

	for i := range parts {
		end := min(start+base, length)

		offsets[i] = []int{start, end - 1}
		start = end
	}
	return offsets
}

func countLines(file *os.File) (int, error) {
	// Create a buffer to read chunks of the file into
	// 32KB is a good default size
	chunkBuffer := make([]byte, 32*1024)
	lineCount := 0
	lineSep := []byte{'\n'}

	for {
		// Read a chunk of the file
		bytesRead, readError := file.Read(chunkBuffer)

		if readError == io.EOF {
			break
		} else if readError != nil {
			return lineCount, readError
		}

		lineCount += bytes.Count(chunkBuffer[:bytesRead], lineSep)
	}
	return lineCount, nil
}
