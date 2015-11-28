package main

import (
	"flag"
	"github.com/FireEater64/gamq"
	log "github.com/cihub/seelog"
	"github.com/davecheney/profile"
	"runtime"
)

func main() {
	// Set up a done channel, that's shared by the whole pipeline.
	// Closing this channel will kill all pipeline goroutines
	//done := make(chan struct{})
	//defer close(done)

	// Set up logging
	initializeLogging()

	// Flush the log before we shutdown
	defer log.Flush()

	// Parse the command line flags
	config := parseCommandLineFlags()
	gamq.SetConfig(&config)

	if config.ProfilingEnabled {
		defer profile.Start(profile.CPUProfile).Stop()
	}

	log.Infof("Broker started on port: %d", gamq.Configuration.Port)
	log.Infof("Executing on: %d threads", runtime.GOMAXPROCS(-1))

	connectionManager := gamq.ConnectionManager{}
	connectionManager.Initialize()
}

func initializeLogging() {
	logger, err := log.LoggerFromConfigAsFile("config/logconfig.xml")

	if err != nil {
		log.Criticalf("An error occurred whilst initializing logging\n", err.Error())
		panic(err)
	}

	log.ReplaceLogger(logger)
}

func parseCommandLineFlags() gamq.Config {
	configToReturn := gamq.Config{}

	flag.IntVar(&configToReturn.Port, "port", 48879, "The port to listen on")
	flag.BoolVar(&configToReturn.ProfilingEnabled, "profile", false, "Produce a pprof file")
	flag.StringVar(&configToReturn.StatsDEndpoint, "statsd", "", "The StatsD endpoint to send metrics to")

	flag.Parse()

	return configToReturn
}
