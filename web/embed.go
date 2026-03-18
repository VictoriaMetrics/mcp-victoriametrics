package web

import "embed"

// DistFS provides the embedded Vite build output.
//
//go:embed dist/*
var DistFS embed.FS
