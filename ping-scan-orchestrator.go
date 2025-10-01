package main

import (
	"io"
	"sync"

	"github.com/cihub/seelog"
)

// --- 3. The Scan Orchestrator ---
// This component brings everything together.

type PingScanOrchestrator struct {
	pinger Pinger
	reader FLSFileReader
	writer IPOutput
	errorTolerance int

	successfulCount int64
	failCount int64
}

// NewPingScanOrchestrator is a factory function that creates
// new orchestrator with its dependencies
func NewPingScanOrchestrator(pinger Pinger, writer IPOutput, reader FLSFileReader, errorTolerance int) *PingScanOrchestrator {
	return &PingScanOrchestrator{
		pinger: pinger,
		writer: writer,
		reader: reader,
		errorTolerance: errorTolerance,
	}
}

const ReadIPErrorMaxThreshold byte = 3

func (o *PingScanOrchestrator) ProcessIPs(count int, outputFile string) {
	var wg sync.WaitGroup

	// A channel to safely collect responsive IPs
	// from concurrent goroutines
	responsiveIPsChan := make(chan string, count)

	var ipsParts [][]int
	if threads > 1 {
		ipsParts = divisionBoundaries(count, threads)
	} else {
		ipsParts = [][]int{{0, count - 1}}
	}

	o.successfulCount = 0
	o.failCount = 0
	for threadIndex := range threads {
		wg.Add(1)

		go o.scanner(threadIndex, ipsParts, responsiveIPsChan, &wg)
	}

	go o.OutputWriter(responsiveIPsChan)
	
	wg.Wait()
	close(responsiveIPsChan)
}

func (o *PingScanOrchestrator) OutputWriter(responsiveIPsChan chan string) {
	// The reachable IP is immediately written
	// into the output file
	var (
		outputFileIsOpen bool = false
		outputFileIO io.ReadWriteCloser
		outputFileError error
		writingBuffer WriterFlusher
	)
	
	if !outputFileIsOpen {
		outputFileIO, outputFileError = o.writer.Open(outputFile)
	
		if outputFileError != nil && !silent {
			seelog.Criticalf(
				"Output write operation failed! Could not write responsive IPs to the output! More: %v",
				outputFileError,
			)
		}
		writingBuffer = o.writer.GetBuffer(outputFileIO)
	}
	for ip := range responsiveIPsChan {
		outputWriteError := o.writer.WriteIP(writingBuffer, ip)
	
		if outputWriteError != nil && !silent {
			seelog.Criticalf(
				"Output write operation failed! Could not write responsive IPs to the output! More: %v",
				outputWriteError,
			)
		}
	}
}

func (o *PingScanOrchestrator) scanner(threadIndex int, ipsParts [][]int, responsiveIPsChan chan string, wg *sync.WaitGroup) {
	inputIPFile, inputIPFileError := o.reader.Open(inputFile)

	if inputIPFileError != nil {
		if !silent {
			seelog.Criticalf(
				"Thread %d: Unable to open input file! More: %s",
				threadIndex,
				inputIPFileError,
			)
		}
		if o.errorTolerance == 0 {
			return
		} else if o.errorTolerance > 0 {
			o.errorTolerance--
		}
	}
	defer inputIPFile.Close()
	
	defer wg.Done()

	var readIPErrors byte = 0
	for lineIndex := ipsParts[threadIndex][0]; lineIndex <= ipsParts[threadIndex][1]; lineIndex++ {
		ipAddr, ipAddrReadError := o.reader.ReadIP(
			int64(lineIndex),
		)

		if ipAddrReadError != nil {
			if !silent {
				seelog.Errorf(
					"Read line operation failure at line %d! More: %s",
					lineIndex,
					ipAddrReadError,
				)
			}
			readIPErrors++

			if readIPErrors >= ReadIPErrorMaxThreshold {
				if o.errorTolerance == 0 {
					return
				} else if o.errorTolerance > 0 {
					o.errorTolerance--
				}
			}
			continue
		}

		if !silent {
			seelog.Infof("Pinging %s", ipAddr)
		}

		pingResult, pingError := o.pinger.Ping(ipAddr)

		if pingError != nil {
			if !silent {
				seelog.Errorf("[Thread %d] %v", threadIndex, pingError)
			}
			if o.errorTolerance == 0 {
				return
			} else if o.errorTolerance > 0 {
				o.errorTolerance--
			}
		}

		if pingResult.Bytes > 0 {
			if !silent {
				seelog.Infof(
					"[Thread %d] %d bytes from %s icmp_seq=%d time=%v",
					threadIndex,
					pingResult.Bytes,
					pingResult.IPAddr,
					pingResult.Sequence,
					pingResult.Latency,
				)
			}
			responsiveIPsChan <- ipAddr

			o.successfulCount++

		} else {
			seelog.Infof("[Thread %d] No response from %s", threadIndex, ipAddr)

			o.failCount++
		}
	}
}
