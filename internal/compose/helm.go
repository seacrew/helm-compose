/*
Copyright © 2023 The Helm Compose Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	cfg "github.com/seacrew/helm-compose/internal/config"
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

func installHelmRelease(name string, release *cfg.Release) {
	var args []string

	args = append(args, "upgrade")
	args = append(args, "--install")

	if release.ChartVersion != "" {
		args = append(args, fmt.Sprintf("--version=%s", release.ChartVersion))
	}

	if release.Namespace != "" {
		args = append(args, fmt.Sprintf("--namespace=%s", release.Namespace))
	}

	if release.ForceUpdate {
		args = append(args, "--force")
	}

	if release.HistoryMax < 0 {
		args = append(args, fmt.Sprintf("--history-max=%d", 0))
	} else if release.HistoryMax > 0 {
		args = append(args, fmt.Sprintf("--history-max=%d", release.HistoryMax))
	}

	if release.CreateNamespace {
		args = append(args, "--create-namespace")
	}

	if release.CleanUpOnFail {
		args = append(args, "--cleanup-on-fail")
	}

	if release.DependencyUpdate {
		args = append(args, "--dependency-update")
	}

	if release.SkipTLSVerify {
		args = append(args, "--insecure-skip-tls-verify")
	}

	if release.SkipCRDs {
		args = append(args, "--skip-crds")
	}

	if release.PostRenderer != "" {
		args = append(args, fmt.Sprintf("--post-renderer=%s", release.PostRenderer))
	}

	if len(release.PostRendererArgs) > 0 {
		args = append(args, fmt.Sprintf("--post-renderer-args=[%s]", strings.Join(release.PostRendererArgs, ",")))
	}

	if release.CAFile != "" {
		args = append(args, fmt.Sprintf("--ca-file=%s", release.CAFile))
	}

	if release.CertFile != "" {
		args = append(args, fmt.Sprintf("--cert-file=%s", release.CertFile))
	}

	if release.KeyFile != "" {
		args = append(args, fmt.Sprintf("--key-file=%s", release.KeyFile))
	}

	if release.Timeout != "" {
		args = append(args, fmt.Sprintf("--timeout=%s", release.Timeout))
	}
	if release.Wait != "" {
		args = append(args, fmt.Sprintf("--wait=%s", release.Wait))
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

	var jsonValues []string
	for key := range release.Values {
		data := util.ConvertJson(release.Values[key])
		values, err := json.Marshal(data)
		if err != nil {
			cp := util.NewColorPrinter(name)
			cp.Printf("%s |\t\t%s", name, err)
			return
		}

		jsonValues = append(jsonValues, fmt.Sprintf("%s=%s", key, values))
	}

	if len(jsonValues) > 0 {
		args = append(args, fmt.Sprintf("--set-json=%s", strings.Join(jsonValues, ",")))
	}

	args = append(args, name)
	args = append(args, release.Chart)

	helmExec(name, args)
}

func uninstallHelmRelease(name string, release *cfg.Release) {
	var args []string

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

	if release.DeletionStrategy != "" {
		args = append(args, fmt.Sprintf("--cascade=%s", release.DeletionStrategy))
	}

	if release.DeletionTimeout != "" {
		args = append(args, fmt.Sprintf("--timeout=%s", release.DeletionTimeout))
	}

	if release.DeletionNoHooks {
		args = append(args, "--no-hooks")
	}

	if release.KeepHistory {
		args = append(args, "--keep-history")
	}

	args = append(args, name)

	helmExec(name, args)
}

func helmExec(name string, args []string) {
	cp := util.NewColorPrinter(name)
	output, _ := util.Execute(helm, args...)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		cp.Printf("%s |\t\t%s", name, scanner.Text())
	}

	err := scanner.Err()

	if err != nil {
		cp.Printf(err.Error())
	}
}
