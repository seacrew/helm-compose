package compose

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/seacrew/helm-compose/internal/util"
)

var (
	helm       = os.Getenv("HELM_BIN")
	versionRE  = regexp.MustCompile(`Version:\s*"([^"]+)"`)
	minVersion = semver.MustParse("v3.0.0")
)

func CompatibleHelmVersion() error {
	cmd := exec.Command(helm, "version")
	util.DebugPrint("Executing %s", strings.Join(cmd.Args, " "))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to run `%s version`: %v", os.Getenv("HELM_BIN"), err)
	}
	versionOutput := string(output)

	matches := versionRE.FindStringSubmatch(versionOutput)
	if matches == nil {
		return fmt.Errorf("Failed to find version in output %#v", versionOutput)
	}
	helmVersion, err := semver.NewVersion(matches[1])
	if err != nil {
		return fmt.Errorf("Failed to parse version %#v: %v", matches[1], err)
	}

	if minVersion.GreaterThan(helmVersion) {
		return fmt.Errorf("helm compose requires at least helm version %s", minVersion.String())
	}
	return nil
}

func addHelmRepository(name string, url string) error {
	output, err := util.Execute(helm, "repo", "add", "--force-update", name, url)

	if err != nil {
		return errors.New(output)
	}

	return nil
}

func installHelmRelease(name string, release *Release) {
	var args []string
	color := util.HashColor(name)

	args = append(args, "upgrade")
	args = append(args, "--install")

	if release.CreateNamespace {
		args = append(args, "--create-namespace")
	}

	if release.ChartVersion != "" {
		args = append(args, fmt.Sprintf("--version=%s", release.ChartVersion))
	}

	if release.Namespace != "" {
		args = append(args, fmt.Sprintf("--namespace=%s", release.Namespace))
	}

	if release.KubeConfig != "" {
		args = append(args, fmt.Sprintf("--kubeconfig=%s", release.KubeConfig))
	}

	if release.KubeContext != "" {
		args = append(args, fmt.Sprintf("--kube-context=%s", release.KubeContext))
	}

	for _, file := range release.ValueFiles {
		args = append(args, fmt.Sprintf("--values=%s", file))
	}

	for _, file := range release.ValueFiles {
		args = append(args, fmt.Sprintf("--values=%s", file))
	}

	var json_values []string
	for key := range release.Values {
		data := util.ConvertJson(release.Values[key])
		values, err := json.Marshal(data)
		if err != nil {
			fmt.Println(color(name + " |\t\tError: " + err.Error()))
			return
		}

		json_values = append(json_values, fmt.Sprintf("%s=%s", key, values))
	}

	if len(json_values) > 0 {
		args = append(args, fmt.Sprintf("--set-json=%s", strings.Join(json_values, ",")))
	}

	args = append(args, name)
	args = append(args, release.Chart)

	output, _ := util.Execute(helm, args...)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		fmt.Println(color(name + " |\t\t" + scanner.Text()))
	}

	err := scanner.Err()

	if err != nil {
		fmt.Print(err.Error())
	}
}

func uninstallHelmRelease(name string, release *Release) {
	var args []string
	color := util.HashColor(name)

	args = append(args, "uninstall")

	if release.Namespace != "" {
		args = append(args, fmt.Sprintf("--namespace=%s", release.Namespace))
	}

	if release.KubeConfig != "" {
		args = append(args, fmt.Sprintf("--kubeconfig=%s", release.KubeConfig))
	}

	if release.KubeContext != "" {
		args = append(args, fmt.Sprintf("--kube-context=%s", release.KubeContext))
	}

	args = append(args, name)

	output, _ := util.Execute(helm, args...)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		fmt.Println(color(name + " |\t\t" + scanner.Text()))
	}

	err := scanner.Err()

	if err != nil {
		fmt.Print(err.Error())
	}
}
