package main

import (
	"flag"
	"log"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/resources"
)

func main() {
	output := flag.String("output", "cmd/mcp-victoriametrics/resources/docsindex/docs-index.tar.gz", "path to the docs search index archive embedded into the server")
	flag.Parse()

	if err := resources.BuildDocsIndex(*output); err != nil {
		log.Fatalf("cannot build docs index: %v", err)
	}
}
