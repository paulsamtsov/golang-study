// Package main provides the Image Metadata Processor service entry point.
package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	processor "github.com/paulsamtsov/lab3-detector/internal/processor"
)

func main() {
	// Setup structured logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Start pprof server on :6060
	go func() {
		log.Info().Msg("Pprof server started on :6060")
		log.Error().Err(http.ListenAndServe("localhost:6060", nil)).Send()
	}()

	// Give pprof time to start
	time.Sleep(100 * time.Millisecond)

	log.Info().Msg("Image Metadata Processor started...")
	log.Info().Msg("Pprof available at http://localhost:6060/debug/pprof/")
	log.Info().Msg("Heap profile: curl http://localhost:6060/debug/pprof/heap > heap.prof")
	log.Info().Msg("CPU profile: curl http://localhost:6060/debug/pprof/profile > cpu.prof")

	// Start worker pool
	processor.RunWorkerPool(5)
}
