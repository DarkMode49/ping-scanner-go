package main

import "sync"

// --- 3. The Scan Orchestrator ---
// This component brings everything together.

type PingScanOrchestrator struct {
	pinger Pinger
	writer IPWriter
}

// NewPingScanOrchestrator is a factory function that creates
// new orchestrator with its dependencies
func NewPingScanOrchestrator(pinger Pinger, writer IPWriter) *PingScanOrchestrator {
	return &PingScanOrchestrator{
		pinger: pinger,
		writer: writer,
	}
}

const ReadIPErrorMaxThreshold byte = 3

func (o *PingScanOrchestrator) ProcessIPs(count int, outputFile string) {
	var wg sync.WaitGroup

	// A channel to safely collect responsive IPs
	// from concurrent goroutines
	responsiveIPsChan := make(chan string, count)

	//TODO int64 version
	var ipsParts [][]int
	if threads > 1 {
		ipsParts = divisionBoundaries(count, threads)
	} else {
		ipsParts = [][]int{{0, count - 1}}
	}

	for threadIndex := range threads {
		wg.Add(1)

		go o.scanner(threadIndex, ipsParts, responsiveIPsChan, &wg)
	}
	
	wg.Wait()
	close(responsiveIPsChan)

	//TODO Rewrite this as a concurrent writing thread
	// Collect all results from the channel into a slice
	var successfulIPs []string
	for ip := range responsiveIPsChan {
		successfulIPs = append(successfulIPs, ip)
	}

	if len(successfulIPs) > 0 {
		outputWriteError := o.writer.WriteIPs(outputFile, successfulIPs)

		if outputWriteError != nil {
			logFatal("Error writing responsive IPs: %v", outputWriteError)
		}
		logMsg(
			"Wrote %d responsive IP addresses to %s",
			len(successfulIPs),
			outputFile,
		)
	} else {
		logMsg("No response from any IP")
	}
}

func (o *PingScanOrchestrator) scanner(threadIndex int, ipsParts [][]int, responsiveIPsChan chan string, wg *sync.WaitGroup) {
	//TODO Repalce with interface to decouple
	//TODO and avoid implementation dependency
	reader := &FileIPReader{}

	inputIPFile, inputIPFileError := reader.Open(inputFile)

	if inputIPFileError != nil {
		logFatal(
			"Thread %d: Unable to open input file! More: %s",
			threadIndex,
			inputIPFileError,
		)
		return
	}
	defer inputIPFile.Close()
	
	defer wg.Done()

	var readIPErrors byte = 0
	for lineIndex := ipsParts[threadIndex][0]; lineIndex <= ipsParts[threadIndex][1]; lineIndex++ {
		ipAddr, ipAddrReadError := reader.ReadIP(
			int64(lineIndex),
		)

		if ipAddrReadError != nil {
			logMsg(
				"Error: Read line operation failure at line %d! More: %s",
				lineIndex,
				ipAddrReadError,
			)
			readIPErrors++

			if readIPErrors >= ReadIPErrorMaxThreshold {
				return
			}
			continue
		}

		logMsg("Pinging %s", ipAddr)

		//TODO IPv6 support
		if o.pinger.Ping(ipAddr) {
			logMsg("SUCCESS: %s", ipAddr)
			responsiveIPsChan <- ipAddr
		} else {
			logMsg("FAILURE: %s", ipAddr)
		}
	}
}
