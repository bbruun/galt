package main

import (
	"fmt"
	"os"

	"github.com/bbruun/galt/common"
	"github.com/bbruun/galt/core"
)

const (
	VERSION string = "0.0.1"
	AUTHOR  string = "Bjarke Bruun"
)

var IS_MASTER bool

func init() {
	/*
		Read command line parameters
	*/

	// Am I a master or minion

}

func main() {
	log := common.ReturnLogger()
	log.Debug("Starting PowerDNS CLI version 0.0.1")

	log.Debugf("%+v", os.Args[1:])
	parameters := []string{}
	if len(os.Args) == 1 {
		log.Debugf("Starting client - using defaults or config file")
		parameters = append(parameters, "client")
	} else {
		log.Debugf("os.Args[0] == %s", os.Args[0])
		for x := 1; x <= len(os.Args[1:]); x++ {
			log.Debugf("- adding os.Args[%s]", os.Args[x])
			parameters = append(parameters, os.Args[x])
		}
	}
	log.Debugf("parameters: %s", parameters[0:])
	if err := core.ParseParameters(parameters); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

}
