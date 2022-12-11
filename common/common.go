package common

import (
	"os"
	"strings"

	"github.com/apsdehal/go-logger"
)

func ReturnLogger() logger.Logger {
	/*
	   Function to return the same log.Logger
	*/
	log, err := logger.New("galt", 1, os.Stdout)
	if err != nil {
		panic(err) // Check for error
	}
	for x := range os.Args {
		if strings.Contains(os.Args[x], "debug") {
			log.SetLogLevel(logger.DebugLevel)
		}
	}
	return *log
}
