package scheduler_docker_local

import (
	"github.com/dokku/dokku/plugins/common"
	"github.com/dokku/dokku/plugins/config"
	"log"
	"os"
)

func TriggerCheckDeploy(appName, containerID, containerType, port, ip, index string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerCorePostDeploy(appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerInstall() error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerPostAppCloneSetup(oldAppName, newAppName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerPostAppRenameSetup(oldAppName, newAppName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerPostCreate(appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerPostDelete(appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerPreDeploy(appName, imageTag string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerPreRestore(scheduler string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerAppStatus(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerDeploy(scheduler, appName, imageTag, processType string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerEnter(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerInspect(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerIsDeployed(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerLogs(scheduler, appName, processType string, tail, prettyPrint bool, num int) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerLogsFailed(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerRegisterRetired(appName, containerID string, wait int) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerRetire(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerRun(scheduler, appName string, envCount int) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerRunList(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

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
