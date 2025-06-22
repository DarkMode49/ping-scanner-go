package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	flag "github.com/spf13/pflag"
)

const version = "0.1.0"
var inputFile = "ips.txt"
var outputFile = "good.txt"
var threads int = 1
const pingTimeout = 2 * time.Second


type IPReader interface {
	ReadIPs(filePath string) ([]string, error)
}

type IPWriter interface {
	WriteIPs(filePath string, ips []string) error
}

type Pinger interface {
	Ping(ipAddress string) bool
}

type FileIPReader struct{}

// ReadIPs reads each line from the given file path and returns it as a slice of strings.
// It adheres to the Single Responsibility Principle (SRP) - its only job is to read.
func (r *FileIPReader) ReadIPs(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open ip file: %w", err)
	}
	defer file.Close()

	var ips []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ips = append(ips, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while scanning file: %w", err)
	}

	return ips, nil
}

// FileIPWriter implements IPWriter to write IPs to a plain text file.
type FileIPWriter struct{}

// WriteIPs writes a slice of IP addresses to the specified file, one IP per line.
// (SRP) - Its only job is to write.
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

// SystemPinger implements the Pinger interface by shelling out to the OS's ping command.
// (SRP) - Its only job is to ping an address.
type SystemPinger struct {
	Timeout time.Duration
}

// Ping executes the system's ping command. It's adaptable for different operating systems.
// This is an example of the Open/Closed Principle. We could create a new LibraryPinger
// that uses a pure Go library instead, without modifying this struct or the code that uses it.
func (p *SystemPinger) Ping(ipAddress string) bool {
	// Create a context with a timeout to prevent
	// pings from hanging indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	var cmd *exec.Cmd

	// The 'ping' command has different arguments on Windows vs. Linux/macOS.
	switch runtime.GOOS {
	case "windows":
		// -n 1: Send 1 echo request.
		// -w: Timeout in milliseconds.
		cmd = exec.CommandContext(
			ctx,
			"ping", "-n", "1", fmt.Sprintf("%d", p.Timeout.Milliseconds()),
			ipAddress,
		)
	default: // Linux, macOS, etc.
		// -c 1: Send 1 packet.
		// -W: Timeout in seconds.
		cmd = exec.CommandContext(
			ctx,
			"ping",
			"-c",
			"1",
			"-W",
			fmt.Sprintf("%.f", p.Timeout.Seconds()),
			ipAddress,
		)
	}

	// We don't care about the output,
	// only whether the command
	// succeeded (exit code 0).
	err := cmd.Run()
	return err == nil
}

// --- 3. The Orchestrator ---
// This component brings everything together.

// PingOrchestrator coordinates the process of reading, pinging, and writing IPs.
// It depends on the interfaces, not the concrete types (Dependency Inversion).
type PingOrchestrator struct {
	reader IPReader
	pinger Pinger
	writer IPWriter
}

// NewPingOrchestrator is a factory function that creates a new orchestrator with its dependencies.
func NewPingOrchestrator(reader IPReader, pinger Pinger, writer IPWriter) *PingOrchestrator {
	return &PingOrchestrator{
		reader: reader,
		pinger: pinger,
		writer: writer,
	}
}

func (o *PingOrchestrator) ProcessIPs(inputFile, outputFile string) {
	var wg sync.WaitGroup
	
	ips, inputFileError := o.reader.ReadIPs(inputFile)
	if inputFileError != nil {
		log.Fatalf("Error: Unable to read IPs: %v", inputFileError)
	}

	log.Printf("Read %d IP addresses.", len(ips))
	log.Println("Starting ICMP ping scan...")

	// A channel to safely collect responsive IPs
	// from concurrent goroutines.
	responsiveIPsChan := make(chan string, len(ips))

	ipsParts := divisionBoundaries(len(ips), threads)
	
	for threadIndex := range threads {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for _, ipAddr := range ips[ipsParts[threadIndex][0]:ipsParts[threadIndex][1]] {
				log.Printf("Pinging %s\n", ipAddr)

				if o.pinger.Ping(ipAddr) {
					log.Printf("SUCCESS: %s is responsive", ipAddr)
					responsiveIPsChan <- ipAddr
				} else {
					log.Printf("FAILURE: %s did not respond", ipAddr)
				}
			}
		}()
	}
	
	wg.Wait()
	close(responsiveIPsChan)

	// Collect all results from the channel into a slice.
	var successfulIPs []string
	for ip := range responsiveIPsChan {
		successfulIPs = append(successfulIPs, ip)
	}

	if len(successfulIPs) > 0 {
		outputWriteError := o.writer.WriteIPs(outputFile, successfulIPs)

		if outputWriteError != nil {
			log.Fatalf("Error writing responsive IPs: %v", outputWriteError)
		}
		log.Printf(
			"Wrote %d responsive IP addresses to %s",
			len(successfulIPs),
			outputFile,
		)
	} else {
		log.Println("No response from any IP")
	}
}

func main() {
	log.Printf("Ping Scanner v%s\n", version)

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
	flag.Parse()

	if threads < 1 {
		fmt.Println("Invalid thread count less than 1!")
		
		os.Exit(1)
	}
	
	if inputFile == "" {
		fmt.Println("Invalid input file parameter!")
		
		os.Exit(1)
	}
	
	if outputFile == "" {
		fmt.Println("Invalid output file parameter!")
		
		os.Exit(1)
	}

	_, ipsFileError := os.Stat(inputFile)

	if errors.Is(ipsFileError, os.ErrNotExist) {
		fmt.Printf("No \"%s\" file was found!\n", inputFile)
		
		os.Exit(1)
	}

	// --- Dependency Injection ---
	// We create our concrete components and "inject" them into the orchestrator.
	// This makes the system modular and testable. You could easily substitute
	// any of these with a different implementation (e.g., a mock for testing).
	reader := &FileIPReader{}
	pinger := &SystemPinger{Timeout: pingTimeout}
	writer := &FileIPWriter{}

	orchestrator := NewPingOrchestrator(reader, pinger, writer)

	orchestrator.ProcessIPs(inputFile, outputFile)
}
