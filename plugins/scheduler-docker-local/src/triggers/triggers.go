package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dokku/dokku/plugins/common"
	"github.com/dokku/dokku/plugins/scheduler-docker-local"
)

// main entrypoint to all triggers
func main() {
	parts := strings.Split(os.Args[0], "/")
	trigger := parts[len(parts)-1]
	flag.Parse()

	var err error
	switch trigger {
	case "check-deploy":
		appName := flag.Arg(0)
		containerID := flag.Arg(1)
		containerType := flag.Arg(2)
		port := flag.Arg(3)
		ip := flag.Arg(4)
		err = scheduler_docker_local.TriggerCheckDeploy(appName, containerID, containerType, port, ip)
	case "install":
		err = scheduler_docker_local.TriggerInstall()
	case "post-create":
		appName := flag.Arg(0)
		err = scheduler_docker_local.TriggerPostCreate(appName)
	case "post-delete":
		appName := flag.Arg(0)
		err = scheduler_docker_local.TriggerPostDelete(appName)
	case "pre-deploy":
		appName := flag.Arg(0)
		imageTag := flag.Arg(1)
		err = scheduler_docker_local.TriggerPreDeploy(appName, imageTag)
	case "pre-restore":
		scheduler := flag.Arg(0)
		err = scheduler_docker_local.TriggerPreRestore(scheduler)
	case "report":
		appName := flag.Arg(0)
		err = scheduler_docker_local.ReportSingleApp(appName, "")
	default:
		common.LogFail(fmt.Sprintf("Invalid plugin trigger call: %s", trigger))
	}

	if err != nil {
		common.LogFail(err.Error())
	}
}
