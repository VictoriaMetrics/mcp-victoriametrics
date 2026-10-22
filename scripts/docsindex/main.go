package main

import (
	"flag"
	"log"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/resources"
)

func main() {
	outputDir := flag.String("outputDir", "cmd/mcp-victoriametrics/resources/docsindex", "directory for the docs search index archives embedded into the server")
	flag.Parse()

	if err := resources.BuildDocsIndexes(*outputDir); err != nil {
		log.Fatalf("cannot build docs indexes: %v", err)
	}
}
