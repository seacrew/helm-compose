package compose

import (
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
	fmt.Printf("Installing release `%s`\n", name)

	var args []string

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

	output, err := util.Execute(helm, args...)

	if err != nil {
		fmt.Print(output)
	}

	fmt.Print(output)
}
