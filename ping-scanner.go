package main

import (
	"errors"
	"io"
	"net"
	"os"
	"time"

	"github.com/cihub/seelog"
	flag "github.com/spf13/pflag"
	"github.com/stoicperlman/fls"
)

const VERSION = "1.0.0"
var inputFile = "ips.txt"
var outputFile = "good.txt"
var threads int = 1
var errorTolerance int = 2
var silent bool = false
const pingTimeout = 2 * time.Second


type PingResult struct {
	Bytes int
	IPAddr *net.IPAddr
	Sequence int
	Latency time.Duration
}

type IPReader interface {
	ReadIP(line int64) (string, error)
}

type FLSFileReader interface {
	Open(filePath string) (*fls.File, error)
	IPReader
}

type IPWriter interface {
	WriteIP(writer WriterFlusher, ip string) error
}

type IPOutput interface {
	Open(filePath string) (io.ReadWriteCloser, error)
	GetBuffer(file io.ReadWriteCloser) (WriterFlusher)
	IPWriter
}

type WriterFlusher interface {
	io.Writer
    Flush() error
}

type Pinger interface {
	Ping(ipAddress string) (PingResult, error)
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
	flag.IntVarP(
		&threads,
		"max-errors",
		"e",
		threads,
		"Number of errors to be ignored (default: 3, -1 for no exit)",
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
		seelog.Error("Invalid thread count less than 1!")
		os.Exit(1)
	}
	
	if inputFile == "" {
		seelog.Error("Invalid input file parameter!")
		os.Exit(1)
	}
	
	if outputFile == "" {
		seelog.Error("Invalid output file parameter!")
		os.Exit(1)
	}
}

func main() {
	arguments()

	loggerInterface, _ := seelog.LoggerFromConfigAsString(seelogConsoleConfig)
	seelog.ReplaceLogger(loggerInterface)

	if !silent {
		seelog.Infof("Ping Scanner v%s", VERSION)
	}

	inputFileStats, ipsFileError := os.Stat(inputFile)

	if errors.Is(ipsFileError, os.ErrNotExist) {
		if !silent {
			seelog.Errorf("No \"%s\" file was found!", inputFile)
		}
		os.Exit(1)
	}

	if inputFileStats.Size() < 8 {
		// 8 bytes is least the file can be having for example:
		// 0.0.0.0
		if !silent {
			seelog.Errorf("Input file \"%s\" is empty or invalid!", inputFile)
		}
		os.Exit(1)
	}

	inputIPFile, inputIPFileError := os.Open(inputFile)
	
	if inputIPFileError != nil {
		if !silent {
			seelog.Criticalf("Unable to open input file! More: %s", inputIPFileError)
		}
		os.Exit(1)
	}
	defer inputIPFile.Close()

	ipCount, ipCountError := countLines(inputIPFile)

	if ipCountError != nil {
		if !silent {
			seelog.Criticalf("Input file loading failure! More: %s", ipCountError)
		}
		return
	}

	pinger := &SystemPinger{Timeout: pingTimeout}
	writer := &FileIPWriter{}
	reader := &FileIPReader{}

	orchestrator := NewPingScanOrchestrator(pinger, writer, reader, errorTolerance)

	if !silent {
		seelog.Warnf("Starting ICMP ping scan...")
	}
	orchestrator.ProcessIPs(ipCount, outputFile)

	if !silent {
		seelog.Info("Scan finished")
		seelog.Infof(
			"Successful: %d  Failed: %d",
			orchestrator.successfulCount,
			orchestrator.failCount,
		)
	}
}
