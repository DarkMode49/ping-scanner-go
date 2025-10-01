package main

import (
	_ "embed"
)

//go:embed seelog-console.xml
var seelogConsoleConfig string

// This will load seelog-console.xml
// into seelogConsoleConfig at runtime
