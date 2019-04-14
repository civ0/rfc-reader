package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Setup logging
	logFile, err := os.OpenFile("rfc-reader.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open logfile: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	_ = UpdateCache()
	index, _ := ReadIndex()
	for i, rfcEntry := range index.RFCEntries {
		fmt.Println(i, ": ", rfcEntry.DocID, " ", rfcEntry.Title, " ", rfcEntry.CanonicalName(), " ", rfcEntry.DocID)
		ReadRFC(rfcEntry.CanonicalName())
	}
}
