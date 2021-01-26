package scheduler_docker_local

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/dokku/dokku/plugins/common"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func TriggerCheckDeploy(appName, containerID, containerType, port, ip string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerCorePostDeploy(appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

// TriggerInstall initializes app restart policies
func TriggerInstall() error {
	if err := common.PropertySetup("scheduler-docker-local"); err != nil {
		return fmt.Errorf("Unable to install the scheduler-docker-local plugin: %s", err.Error())
	}

	// mkdir -p "${DOKKU_LIB_ROOT}/data/scheduler-docker-local"
	directory := filepath.Join(common.MustGetEnv("DOKKU_LIB_ROOT"), "data", "scheduler-docker-local")
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}
	// chown -R "${DOKKU_SYSTEM_USER}:${DOKKU_SYSTEM_GROUP}" "${DOKKU_LIB_ROOT}/data/scheduler-docker-local"
	systemUser := common.GetenvWithDefault("DOKKU_SYSTEM_USER", "dokku")
	systemGroup := common.GetenvWithDefault("DOKKU_SYSTEM_GROUP", "dokku")
	err := filepath.Walk(directory, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return chown(systemUser, systemGroup, directory)
	})
	if err != nil {
		return err
	}

	// fn-plugin-property-setup "scheduler-docker-local"
	if err := common.PropertySetup("scheduler-docker-local"); err != nil {
		return err
	}

	// DOKKU_PATH="$(which dokku)"
	dokkuPath, err := exec.LookPath("dokku")
	if err != nil {
		return err
	}
	found := false
	b, err := common.NewShellCmdWithArgs("systemctl", "--no-pager").Output()
	if err != nil {
		return err
	}
	s := bufio.NewScanner(bytes.NewReader(b))
	for s.Scan() {
		l := s.Text()
		if strings.Contains(l, "-.mount") {
			found = true
			break
		}
	}
	// if [[ $(systemctl 2>/dev/null) =~ -\.mount ]]; then
	if found {
		// TODO: consider using template
		f, err := os.Create("/etc/systemd/system/dokku-retire.service")
		if err != nil {
			return err
		}
		_, err = f.WriteString(`[Unit]
Description=Dokku retire service
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
User=` + systemUser + `
ExecStart=` + dokkuPath + ` ps:retire

[Install]
WantedBy=docker.service
`)
		if err != nil {
			return err
		}
		f, err = os.Create("/etc/systemd/system/dokku-retire.timer")
		if err != nil {
			return err
		}
		_, err = f.WriteString(`[Unit]
Description=Run dokku-retire.service every 5 minutes

[Timer]
OnCalendar=*:0/5
Persistent=true

[Install]
WantedBy=timers.target
`)
		if err != nil {
			return err
		}
		// if command -v systemctl &>/dev/null; then
		if common.NewShellCmd("command -v systemctl &>/dev/null").Execute() {
			common.NewShellCmd("systemctl --quiet reenable dokku-retire").Execute()
			common.NewShellCmd("systemctl --quiet enable dokku-retire.timer").Execute()
			common.NewShellCmd("systemctl --quiet start dokku-retire.timer").Execute()
		}
	} else {
		f, err := os.Create("/etc/cron.d/dokku-retire")
		if err != nil {
			return err
		}
		// TODO: consider using template
		_, err = f.WriteString(`PATH=/usr/local/bin:/usr/bin:/bin
SHELL=/bin/bash

*/5 * * * * ` + systemUser + ` ` + dokkuPath + ` ps:retire >> /var/log/dokku/retire.log 2>&1`)
	}
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
	p := filepath.Join(common.MustGetEnv("DOKKU_LIB_ROOT"), "data", "scheduler-docker-local", appName)
	// mkdir -p "${DOKKU_LIB_ROOT}/data/scheduler-docker-local/${APP}"
	if err := os.MkdirAll(p, 0755); err != nil {
		return err
	}
	systemUser := common.GetenvWithDefault("DOKKU_SYSTEM_USER", "dokku")
	systemGroup := common.GetenvWithDefault("DOKKU_SYSTEM_GROUP", "dokku")
	// chown -R "${DOKKU_SYSTEM_USER}:${DOKKU_SYSTEM_GROUP}" "${DOKKU_LIB_ROOT}/data/scheduler-docker-local/${APP}"
	return chown(systemUser, systemGroup, p)
}

func TriggerPostDelete(appName string) error {
	// fn-plugin-property-destroy "scheduler-docker-local" "$APP"
	if err := common.PropertyDestroy("scheduler-docker-local", appName); err != nil {
		return err
	}
	// rm -rf "${DOKKU_LIB_ROOT}/data/scheduler-docker-local/$APP"
	if err := os.RemoveAll(filepath.Join(common.MustGetEnv("DOKKU_LIB_ROOT"), "data", "scheduler-docker-local", appName)); err != nil {
		return err
	}

	// local DOKKU_SCHEDULER=$(get_app_scheduler "$APP")
	scheduler := common.GetAppScheduler(appName)
	// if [[ "$DOKKU_SCHEDULER" != "docker-local" ]]; then
	if scheduler != "docker-local" {
		return nil
	}

	// local IMAGE_REPO=$(get_app_image_repo "$APP")
	repo := common.GetAppImageRepo(appName)

	// # remove all application containers & images
	// # shellcheck disable=SC2046
	// local DOKKU_APP_CIDS=$("$DOCKER_BIN" container list --all --no-trunc | grep "dokku/${APP}:" | awk '{ print $1 }' | xargs)
	b, err := common.NewShellCmdWithArgs(
		common.DockerBin(),
		"container", "list", "--all", "--no-trunc", "--format", "{{.ID}},{{.Image}}",
	).Output()
	if err != nil {
		return err
	}
	s := bufio.NewScanner(bytes.NewReader(b))
	cnamePrefix := fmt.Sprintf("dokku/%s:", appName)
	// if [[ -n "$DOKKU_APP_CIDS" ]]; then
	var ids []string
	for s.Scan() {
		l := s.Text()
		if !strings.Contains(l, cnamePrefix) {
			continue
		}
		ids = append(ids, strings.SplitN(l, ",", 2)[1])
	}
	// "$DOCKER_BIN" container rm --force $DOKKU_APP_CIDS >/dev/null 2>&1 || true
	args := []string{"container", "rm", "--force"}
	args = append(args, ids...)
	cmd := common.NewShellCmdWithArgs(common.DockerBin(), args...)
	cmd.ShowOutput = false
	cmd.Execute()

	// "$DOCKER_BIN" image rm $("$DOCKER_BIN" image list --quiet "$IMAGE_REPO" | xargs) &>/dev/null || true
	cmd = common.NewShellCmdWithArgs(common.DockerBin(), "image", "list", "--quiet", repo)
	cmd.ShowOutput = false
	b, err = cmd.Output()
	if err != nil {
		return err
	}
	s = bufio.NewScanner(bytes.NewReader(b))
	ids = nil
	for s.Scan() {
		l := s.Text()
		ids = append(ids, l)
	}
	args = []string{"image", "rm"}
	args = append(args, ids...)
	cmd = common.NewShellCmdWithArgs(common.DockerBin(), args...)
	cmd.ShowOutput = false
	cmd.Execute()
	return nil
}

func TriggerPreDeploy(appName, imageTag string) error {
	// local DOKKU_SCHEDULER=$(get_app_scheduler "$APP")
	scheduler := common.GetAppScheduler(appName)
	// if [[ "$DOKKU_SCHEDULER" != "docker-local" ]]; then
	if scheduler != "docker-local" {
		return nil
	}

	// scheduler-docker-local-pre-deploy-chown-app "$APP" "$IMAGE_TAG"
	if err := chownApp(appName, imageTag); err != nil {
		return err
	}
	// scheduler-docker-local-pre-deploy-precheck "$APP" "$IMAGE_TAG"
	if err := preCheck(appName, imageTag, "TriggerPreDeploy"); err != nil {
		return err
	}
	return nil
}

func TriggerPreRestore(scheduler string) error {
	// if [[ "$DOKKU_SCHEDULER" != "docker-local" ]]; then
	if scheduler == "docker-local" {
		return nil
	}

	// "$DOCKER_BIN" container rm $("$DOCKER_BIN" container list --all --format "{{.Names}}" --filter "label=$DOKKU_CONTAINER_LABEL" --quiet | grep -E '(.+\..+\.[0-9]+\.[0-9]+$)') &>/dev/null || true
	filter := fmt.Sprintf("label=%v", os.Getenv("DOKKU_CONTAINER_LABEL"))
	out, err := common.NewShellCmdWithArgs(common.DockerBin(), "container", "list", "--all", "--format", "{{.Names}}", "--filter", filter, "--quiet").Output()
	if err != nil {
		return err
	}
	pattern := regexp.MustCompilePOSIX(`(.+\..+\.[0-9]+\.[0-9]+$)`)
	s := bufio.NewScanner(bytes.NewReader(out))
	var containers []string
	for s.Scan() {
		l := s.Text()
		if pattern.MatchString(l) {
			containers = append(containers, l)
		}
	}
	var args []string
	args = append(args, "container", "rm")
	args = append(args, containers...)
	cmd := common.NewShellCmdWithArgs(common.DockerBin(), args...)
	cmd.ShowOutput = false
	cmd.Execute()
	return nil
}

func TriggerSchedulerAppStatus(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerDeploy(scheduler, appName, imageTag string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerDockerCleanup(scheduler, appName string, force bool) error {
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

func TriggerSchedulerRegisterRetired(scheduler, appName string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerRetire() error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerRun(scheduler, appName string, envCount int) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerStop(scheduler, appName string, removeContainers bool) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerTagsCreate(scheduler, appName, sourceImage, targetImage string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func TriggerSchedulerTagsDestroy(scheduler, appName, imageRepo, imageTag string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}
