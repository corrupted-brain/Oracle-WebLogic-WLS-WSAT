package libcve201710271

import (
	"fmt"
	"strings"
)

// Banner prints out a banner with the settings of the applications
func Banner(config Config) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("Author: Kevin Kirsche (d3c3pt10n)")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("Configuration:")
	fmt.Printf("\tListening Host: %s\n", config.Lhost)
	fmt.Printf("\tListening Port: %d\n", config.Lport)
	fmt.Printf("\tOutput File: %s\n", config.OutputFile)
	fmt.Printf("\tTargets File: %s\n", config.TargetFile)
	fmt.Printf("\tThreads: %d\n", config.Threads)
	fmt.Printf("\tScan Complete Wait Time: %d\n", config.WaitTime)
	fmt.Printf("\tScan All URLs: %t\n", config.AllURLs)
	fmt.Printf("\tVerbose mode: %t\n", config.Verbose)
	fmt.Println(strings.Repeat("=", 80))
}
