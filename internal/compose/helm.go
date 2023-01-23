package compose

import (
	"bufio"
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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

var colors [](func(...interface{}) string) = [](func(...interface{}) string){
	Color("%s"), // fallback
	Color("\033[1;32m%s\033[0m"),
	Color("\033[1;33m%s\033[0m"),
	Color("\033[1;34m%s\033[0m"),
	Color("\033[1;35m%s\033[0m"),
	Color("\033[1;36m%s\033[0m"),
	Color("\033[1;90m%s\033[0m"),
	Color("\033[1;92m%s\033[0m"),
	Color("\033[1;93m%s\033[0m"),
	Color("\033[1;94m%s\033[0m"),
	Color("\033[1;95m%s\033[0m"),
	Color("\033[1;96m%s\033[0m"),
}

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func hashColor(s string) func(...interface{}) string {
	h := fnv.New32a()
	h.Write([]byte(s))

	hash := fmt.Sprint(h.Sum32())
	for {
		subtotal := 0
		for _, r := range hash {
			value, _ := strconv.Atoi(string(r))
			subtotal += value
		}

		if subtotal < len(colors) {
			return colors[subtotal]
		}

		hash = fmt.Sprint(subtotal)
	}
}

func installHelmRelease(name string, release *Release) {
	var args []string
	color := hashColor(name)

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
	color := hashColor(name)

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
