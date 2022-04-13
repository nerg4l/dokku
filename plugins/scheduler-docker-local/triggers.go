package scheduler_docker_local

import (
	"fmt"
	"github.com/dokku/dokku/plugins/common"
	"github.com/dokku/dokku/plugins/config"
	"log"
	"os"
)

// TriggerCheckDeploy scheduler-docker-local check-deploy plugin trigger
func TriggerCheckDeploy(appName, containerID, containerType, port, ip, index string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerCorePostDeploy scheduler-docker-local core-post-deploy state cleanup
func TriggerCorePostDeploy(appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerInstall scheduler-docker-local install plugin trigger
func TriggerInstall() error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerPostAppCloneSetup removes docker-local files when setting up a clone
func TriggerPostAppCloneSetup(oldAppName, newAppName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerPostAppRenameSetup updates settings when renaming an app
func TriggerPostAppRenameSetup(oldAppName, newAppName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerPostCreate scheduler-docker-local post-create plugin trigger
func TriggerPostCreate(appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerPostDelete scheduler-docker-local post-delete plugin trigger
func TriggerPostDelete(appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerPreDeploy scheduler-docker-local pre-deploy plugin trigger
func TriggerPreDeploy(appName, imageTag string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerPreRestore scheduler-docker-local pre-restore plugin trigger
func TriggerPreRestore(scheduler string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerSchedulerAppStatus fetches the status for a given app
func TriggerSchedulerAppStatus(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerSchedulerDeploy deploys an image tag for a given application
func TriggerSchedulerDeploy(scheduler, appName, imageTag, processType string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerSchedulerEnter enter a running container
func TriggerSchedulerEnter(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerSchedulerInspect scheduler-docker-local scheduler-inspect plugin trigger
func TriggerSchedulerInspect(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerSchedulerIsDeployed checks if an app is deployed
func TriggerSchedulerIsDeployed(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerSchedulerLogs cheduler-docker-local scheduler-logs plugin trigger
func TriggerSchedulerLogs(scheduler, appName, processType string, tail, prettyPrint bool, num int) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerSchedulerLogsFailed scheduler-docker-local scheduler-logs-failed plugin trigger
func TriggerSchedulerLogsFailed(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerSchedulerRegisterRetired register a container for retiring
func TriggerSchedulerRegisterRetired(appName, containerID string, wait int) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerSchedulerRetire retires all old containers once they have aged out
func TriggerSchedulerRetire(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerSchedulerRun runs command in container based on app image
func TriggerSchedulerRun(scheduler, appName string, runEnv ...string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerSchedulerRunList runs command in container based on app image
func TriggerSchedulerRunList(scheduler, appName string) error {
	if err := common.VerifyAppName(appName); err != nil {
		return err
	}
	if scheduler != "docker-local" {
		return nil
	}

	common.LogInfo2Quiet(fmt.Sprintf("%s run containers", appName))
	cmd := common.NewShellCmdWithArgs(
		common.DockerBin(), "ps",
		"--filter", fmt.Sprintf("label=com.dokku.app-name=%s", appName),
		"--filter", "label=com.dokku.container-type=run",
		"--format", "table {{.Names}}\t{{.Command}}\t{{.RunningFor}}",
	)
	cmd.Execute()
	return nil
}

// TriggerSchedulerStop scheduler-docker-local scheduler-stop plugin trigger
func TriggerSchedulerStop(scheduler, appName string, removeContainers bool) error {
	if scheduler != "docker-local" {
		return nil
	}

	containerIDs, err := common.GetAppRunningContainerIDs(appName, "")
	if err != nil {
		return err
	}
	stopTimeout, _ := config.Get(appName, "DOKKU_DOCKER_STOP_TIMEOUT")

	if len(containerIDs) > 0 {
		// Disable the container restart policy
		args := []string{"container", "update", "--restart", "no"}
		args = append(args, containerIDs...)
		cmd := common.NewShellCmdWithArgs(common.DockerBin(), args...)
		cmd.ShowOutput = false
		cmd.Command.Stderr = os.Stderr
		cmd.Execute()

		args = []string{"container", "stop"}
		if stopTimeout != "" {
			args = append(args, "--time", stopTimeout)
		}
		args = append(args, containerIDs...)
		cmd = common.NewShellCmdWithArgs(common.DockerBin(), args...)
		cmd.ShowOutput = false
		cmd.Command.Stderr = os.Stderr
		cmd.Execute()
	}

	if removeContainers {
		cids, err := common.GetAppContainerIDs(appName, "")
		if err != nil {
			return err
		}
		if len(cids) > 0 {
			for i := range cids {
				if err := common.PlugnTrigger("scheduler-register-retired", appName, cids[i]); err != nil {
					return err
				}
			}
			args := []string{"container", "rm", "--force"}
			args = append(args, cids...)
			cmd := common.NewShellCmdWithArgs(common.DockerBin(), args...)
			cmd.ShowOutput = false
			cmd.Execute()
		}
	}
	return nil
}
