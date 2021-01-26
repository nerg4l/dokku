package main

import (
	"fmt"
	scheduler_docker_local "github.com/dokku/dokku/plugins/scheduler-docker-local"
	"os"
	"strings"

	"github.com/dokku/dokku/plugins/common"

	flag "github.com/spf13/pflag"
)

// main entrypoint to all subcommands
func main() {
	parts := strings.Split(os.Args[0], "/")
	subcommand := parts[len(parts)-1]

	var err error
	switch subcommand {
	case "report":
		args := flag.NewFlagSet("scheduler-docker-local:report", flag.ExitOnError)
		osArgs, infoFlag, flagErr := common.ParseReportArgs("scheduler-docker-local", os.Args[2:])
		if flagErr == nil {
			args.Parse(osArgs)
			appName := args.Arg(0)
			err = scheduler_docker_local.CommandReport(appName, infoFlag)
		}
	case "set":
		args := flag.NewFlagSet("scheduler-docker-local:set", flag.ExitOnError)
		args.Parse(os.Args[2:])
		appName := args.Arg(0)
		property := args.Arg(1)
		value := args.Arg(2)
		err = scheduler_docker_local.CommandSet(appName, property, value)
	default:
		common.LogFail(fmt.Sprintf("Invalid plugin subcommand call: %s", subcommand))
	}

	if err != nil {
		common.LogFail(err.Error())
	}
}
