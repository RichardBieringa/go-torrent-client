package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RichardBieringa/go-torrent-client/torrentfile"
)

func printUsage() {
	fmt.Printf("Usage: %v torrentfile destination\n", os.Args[0])
}

func main() {
	log.Println("Starting Go Torrent Client!")

	// Usage Check
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	// Get torrent file from command line arg
	torrentFilePath := os.Args[1]

	// Get destination directory from command line arg
	desinationPath := os.Args[2]

	// debug
	fmt.Printf("Torrent File: %v\n", torrentFilePath)
	fmt.Printf("Destination Directory: %v\n", desinationPath)

	torrentfile.Open(torrentFilePath)

}
