package core

import (
	"fmt"

	"github.com/bbruun/galt/client"
	"github.com/bbruun/galt/common"
	"github.com/bbruun/galt/server"
)

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func ParseParameters(args []string) error {
	log := common.ReturnLogger()
	subcommand := FilterSubCommand(args[0])
	log.Debugf("subcommand: %s", subcommand)

	// start the proper backend
	cmds := []Runner{}
	if subcommand == "server" {
		cmds = []Runner{
			server.NewServer(),
		}
	} else {
		// Default to "client" mode
		subcommand = "client"
	}
	if subcommand == "client" {
		log.Debugf("running: client")
		cmds = []Runner{
			client.NewClient(),
		}
	}

	// add the rest of the parameters to the backend
	for _, cmd := range cmds {
		log.Debugf("- %+v", cmd)
		if cmd.Name() == subcommand {
			log.Debugf("- running %s", subcommand)
			cmd.Init(args)
			return cmd.Run()
		} else {
			log.Debugf("%s is not runnable", cmd)
		}
	}

	return fmt.Errorf("unknown parameter passed")
}

func FilterSubCommand(subcommand string) string {
	log := common.ReturnLogger()
	// Valid subcommands that map to the real subcommand that is returned
	log.Debugf("filter command: %s", subcommand)
	masterOrClient := map[string]bool{
		"server": true,
	}
	if masterOrClient[subcommand] {
		log.Debugf("parameter \"server\" provided, running as server")
		return "server"
	}

	log.Debugf("no parameter provided, running as client")
	// Default run as client
	return "client"
}
