package main

import (
	"fmt"
	"os"
	"time"
)

const LogSampleTimeFormat = "[2006-01-02 15:04:05]"

func logMsg(format string, a ...any) {
	if silent {
		return
	}
	now := time.Now()
	timestamp := now.Format(LogSampleTimeFormat)

	fmt.Printf("%s %s\n", timestamp, fmt.Sprintf(format, a...))
}

func logFatal(format string, a ...any) {
	if silent {
		return
	}
	now := time.Now()
	timestamp := now.Format(LogSampleTimeFormat)

	fmt.Fprintf(os.Stderr, "%s %s\n", timestamp, fmt.Sprintf(format, a...))
}
