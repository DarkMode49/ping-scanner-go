package main

import (
	"errors"
	"io"
	"log"
	"os"
	"time"

	flag "github.com/spf13/pflag"
)

const version = "0.2.0"
var inputFile = "ips.txt"
var outputFile = "good.txt"
var threads int = 1
var silent bool = false
const pingTimeout = 2 * time.Second


type IPReader interface {
	ReadIP() (string, error)
}

type IPWriter interface {
	WriteIPs(filePath string, ips []string) error
}

type Pinger interface {
	Ping(ipAddress string) bool
}


func arguments() {
	flag.StringVarP(
		&inputFile,
		"input",
		"i",
		inputFile,
		"File with IP addresses in consequent lines",
	)
	flag.StringVarP(
		&outputFile,
		"output",
		"o",
		outputFile,
		"File with IP addresses in consequent lines",
	)
	flag.IntVarP(
		&threads,
		"threads",
		"h",
		threads,
		"Number of concurrent threads (default: 1, recommended: number of CPU cores)",
	)
	flag.BoolVarP(
		&silent,
		"silent",
		"s",
		silent,
		"No console print",
	)
	flag.Parse()
	
	if threads < 1 {
		logMsg("Invalid thread count less than 1!")
		os.Exit(1)
	}
	
	if inputFile == "" {
		logMsg("Invalid input file parameter!")
		os.Exit(1)
	}
	
	if outputFile == "" {
		logMsg("Invalid output file parameter!")
		os.Exit(1)
	}
}

//TODO Improve logging by adding proper level
func main() {
	log.SetFlags(0)
	log.SetOutput(io.Discard)

	logMsg("Ping Scanner v%s", version)

	arguments()

	inputFileStats, ipsFileError := os.Stat(inputFile)

	if errors.Is(ipsFileError, os.ErrNotExist) {
		logMsg("No \"%s\" file was found!", inputFile)
		os.Exit(1)
	}

	if inputFileStats.Size() < 8 {
		// 8 bytes is least the file can be having for example:
		// 0.0.0.0
		logMsg("Input file \"%s\" is empty or invalid!", inputFile)
		os.Exit(1)
	}

	// Dependency Injection
	inputIPFile, inputIPFileError := os.Open(inputFile)
	
	if inputIPFileError != nil {
		logFatal("Unable to open input file! More: %s", inputIPFileError)
		os.Exit(1)
	}
	defer inputIPFile.Close()

	ipCount, ipCountError := countLines(inputIPFile)

	if ipCountError != nil {
		logFatal("Error: Input file loading failure! More: %s", ipCountError)
		os.Exit(1)
	}

	pinger := &SystemPinger{Timeout: pingTimeout}
	writer := &FileIPWriter{}

	orchestrator := NewPingScanOrchestrator(pinger, writer)

	logMsg("Starting ICMP ping scan...")
	orchestrator.ProcessIPs(ipCount, outputFile)
}
