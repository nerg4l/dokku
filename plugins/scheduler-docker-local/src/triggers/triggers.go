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
		index := flag.Arg(5)
		err = scheduler_docker_local.TriggerCheckDeploy(appName, containerID, containerType, port, ip, index)
	case "core-post-deploy":
		appName := flag.Arg(0)
		err = scheduler_docker_local.TriggerCorePostDeploy(appName)
	case "install":
		err = scheduler_docker_local.TriggerInstall()
	case "post-app-clone-setup":
		oldAppName := flag.Arg(0)
		newAppName := flag.Arg(1)
		err = scheduler_docker_local.TriggerPostAppCloneSetup(oldAppName, newAppName)
	case "post-app-rename-setup":
		oldAppName := flag.Arg(0)
		newAppName := flag.Arg(1)
		err = scheduler_docker_local.TriggerPostAppRenameSetup(oldAppName, newAppName)
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
	case "scheduler-app-status":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		err = scheduler_docker_local.TriggerSchedulerAppStatus(scheduler, appName)
	case "scheduler-deploy":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		imageTag := flag.Arg(2)
		processType := flag.Arg(3)
		err = scheduler_docker_local.TriggerSchedulerDeploy(scheduler, appName, imageTag, processType)
	case "scheduler-enter":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		err = scheduler_docker_local.TriggerSchedulerEnter(scheduler, appName)
	case "scheduler-inspect":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		err = scheduler_docker_local.TriggerSchedulerInspect(scheduler, appName)
	case "scheduler-is-deployed":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		err = scheduler_docker_local.TriggerSchedulerIsDeployed(scheduler, appName)
	case "scheduler-logs":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		processType := flag.Arg(2)
		tail := common.ToBool(flag.Arg(3))
		pretty := common.ToBool(flag.Arg(4))
		// TODO: check default value
		num := common.ToInt(flag.Arg(5), 0)
		err = scheduler_docker_local.TriggerSchedulerLogs(scheduler, appName, processType, tail, pretty, num)
	case "scheduler-logs-failed":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		err = scheduler_docker_local.TriggerSchedulerLogsFailed(scheduler, appName)
	case "scheduler-register-retired":
		appName := flag.Arg(0)
		containerID := flag.Arg(1)
		// TODO: check type
		wait := common.ToInt(flag.Arg(2), -60)
		err = scheduler_docker_local.TriggerSchedulerRegisterRetired(appName, containerID, wait)
	case "scheduler-retire":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		err = scheduler_docker_local.TriggerSchedulerRetire(scheduler, appName)
	case "scheduler-run":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		// TODO: check default value
		envCount := common.ToInt(flag.Arg(2), 0)
		err = scheduler_docker_local.TriggerSchedulerRun(scheduler, appName, flag.Args()[3:3+envCount]...)
	case "scheduler-run-list":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		err = scheduler_docker_local.TriggerSchedulerRunList(scheduler, appName)
	case "scheduler-stop":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		removeContainers := common.ToBool(flag.Arg(1))
		err = scheduler_docker_local.TriggerSchedulerStop(scheduler, appName, removeContainers)
	default:
		common.LogFail(fmt.Sprintf("Invalid plugin trigger call: %s", trigger))
	}

	if err != nil {
		common.LogFail(err.Error())
	}
}
