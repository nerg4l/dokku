package scheduler_docker_local

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/dokku/dokku/plugins/common"
	"github.com/dokku/dokku/plugins/config"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
)

func RetireContainer(appName, cid, deadTime string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func RetireContainers() error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func RetireImages() error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func RegisterRetired(imageType, appName, dockerID, wait string) error {
	// TODO: implement
	log.Fatal("not implemented")
	return nil
}

func chown(username, groupname, p string) error {
	u, err := user.Lookup(username)
	if err != nil {
		return err
	}
	g, err := user.LookupGroup(groupname)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return err
	}
	return os.Chown(p, uid, gid)
}

func chownApp(appName, imageTag string) error {
	// local DOCKER_RUN_LABEL_ARGS="--label=com.dokku.app-name=$APP"
	dockerRunLabelArgs := []string{fmt.Sprintf("--label=com.dokku.app-name=%s", appName)}

	// IMAGE=$(get_app_image_name "$APP" "$IMAGE_TAG")
	image := common.GetAppImageName(appName, imageTag, "")

	// DOKKU_APP_TYPE=$(config_get "$APP" DOKKU_APP_TYPE || true)
	appType, _ := config.Get(appName, "DOKKU_APP_TYPE")
	// DOKKU_APP_USER=$(config_get "$APP" DOKKU_APP_USER || true)
	appUser, _ := config.Get(appName, "DOKKU_APP_USER")
	if appUser == "" {
		// DOKKU_APP_USER=${DOKKU_APP_USER:="herokuishuser"}
		appUser = "herokuishuser"
	}
	// APP_PATHS=$(plugn trigger storage-list "$APP" "deploy")
	b, err := common.PlugnTriggerOutput("storage-list", appName, "deploy")
	if err != nil {
		return err
	}

	var containerPaths []string
	var dockerArgs []string
	// if [[ -n "$APP_PATHS" ]]; then
	if len(b) > 0 {
		// CONTAINER_PATHS=$(echo "$APP_PATHS" | awk -F ':' '{ print $2 }' | xargs)
		s := bufio.NewScanner(bytes.NewReader(b))
		for s.Scan() {
			containerPaths = append(containerPaths, strings.SplitN(s.Text(), ":", 2)[1])
		}
		// TODO: DOCKER_ARGS=$(: | plugn trigger docker-args-deploy "$APP" "$IMAGE_TAG")
		//  What does `: |` do?
		bb, err := common.PlugnTriggerOutput("docker-args-deploy", appName, imageTag)
		if err != nil {
			return err
		}
		// # strip --restart args from DOCKER_ARGS
		// DOCKER_ARGS=$(sed -e "s/--restart=[[:graph:]]\+[[:blank:]]\?//g" <<<"$DOCKER_ARGS")
		// TODO: Check `sed` command output.
		pattern := regexp.MustCompile(`--restart=[[:graph:]]\+[[:blank:]]\?`)
		dockerArgs = strings.Split(pattern.ReplaceAllString(string(bb), ""), " ")
	}

	// if [[ "$DOKKU_APP_TYPE" != "herokuish" ]] || [[ -z "$CONTAINER_PATHS" ]]; then
	if appType != "herokuish" || len(containerPaths) == 0 {
		return nil
	}

	// DISABLE_CHOWN="$(fn-plugin-property-get "scheduler-docker-local" "$APP" "disable-chown" "")"
	p := common.PropertyGet("scheduler-docker-local", appName, "disable-chown")
	// if [[ "$DISABLE_CHOWN" == "true" ]]; then
	if p == "true" {
		return nil
	}

	// # shellcheck disable=SC2086
	// "$DOCKER_BIN" container run --rm "${DOCKER_RUN_LABEL_ARGS[@]}" $DOKKU_GLOBAL_RUN_ARGS "${ARG_ARRAY[@]}" $IMAGE /bin/bash -c "find $CONTAINER_PATHS -not -user $DOKKU_APP_USER -print0 | xargs -0 -r chown -R $DOKKU_APP_USER" || true
	var args []string
	args = append(args, "container", "run", "--rm")
	args = append(args, dockerRunLabelArgs...)
	args = append(args, strings.Split(common.MustGetEnv("DOKKU_GLOBAL_RUN_ARGS"), " ")...)
	args = append(args, dockerArgs...)
	// TODO: break down commands
	findArg := fmt.Sprintf("find %s -not -user %s -print0 | xargs -0 -r chown -R %s", containerPaths, appUser, appUser)
	args = append(args, image, "/bin/bash", "-c", findArg)
	common.NewShellCmdWithArgs(common.DockerBin(), args...).Execute()
	log.Fatal("not implemented")
	return nil
}

func preCheck(appName, imageTag, funcName string) error {
	// local IMAGE=$(get_deploying_app_image_name "$APP" "$IMAGE_TAG")
	img, err := common.GetDeployingAppImageName(appName, imageTag, "")
	if err != nil {
		return err
	}
	// local CHECKS_FILE=$(mktemp "/tmp/dokku-${DOKKU_PID}-${FUNCNAME[0]}.XXXXXX")
	f, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("dokku-%d-%s.*", os.Getpid(), funcName))
	if err != nil {
		return err
	}
	f.Close() // Release file to allow modification
	// copy_from_image "$IMAGE" "CHECKS" "$CHECKS_FILE" 2>/dev/null || true
	_ = common.CopyFromImage(appName, img, "CHECKS", f.Name())

	// dokku_log_info2 "Processing deployment checks"
	common.LogInfo2("Processing deployment checks")
	stat, err := os.Stat(f.Name())
	if err != nil {
		return err
	}
	// if [[ ! -s "${CHECKS_FILE}" ]]; then
	if stat.Size() == 0 {
		//   local CHECKS_URL="${DOKKU_CHECKS_URL:-http://dokku.viewdocs.io/dokku/deployment/zero-downtime-deploys/}"
		checksUrl := common.GetenvWithDefault("DOKKU_CHECKS_URL", "http://dokku.viewdocs.io/dokku/deployment/zero-downtime-deploys/")
		//   dokku_log_verbose "No CHECKS file found. Simple container checks will be performed."
		common.LogVerbose("No CHECKS file found. Simple container checks will be performed.")
		//   dokku_log_verbose "For more efficient zero downtime deployments, create a CHECKS file. See ${CHECKS_URL} for examples"
		common.LogVerbose(fmt.Sprintf("For more efficient zero downtime deployments, create a CHECKS file. See %s for examples", checksUrl))
	}
	return nil
}
