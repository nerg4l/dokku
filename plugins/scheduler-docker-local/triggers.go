package scheduler_docker_local

import (
	"log"
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
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}
