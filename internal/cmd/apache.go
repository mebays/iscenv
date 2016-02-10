/*
Copyright 2016 Ontario Systems

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

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ontariosystems/iscenv/internal/iscenv"

	"github.com/spf13/cobra"
)

const (
	gatewayDir    = "/opt/iscenv-csp-gateway"
	cspIniPath    = gatewayDir + "/bin/CSP.ini"
	srcGatewayDir = "/ensemble/cspgateway/."
	apacheDir     = "/etc/apache2"
)

var apacheCmd = &cobra.Command{
	Use:   "apache INSTANCE [INSTANCE...]",
	Short: "Create an Apache site for the instances",
	Long:  "Attempt to install/configure the CSP gateway and Apache.  This command must be run as root.",
	Run:   configureApacheSite,
}

func init() {
	rootCmd.AddCommand(apacheCmd)

	addMultiInstanceFlags(apacheCmd, "apache")
}

func configureApacheSite(_ *cobra.Command, args []string) {
	iscenv.Ensure(iscenv.IsUserRoot)

	instances := multiInstanceFlags.getInstances(args)
	var validInstances iscenv.ISCInstances
	for _, instanceName := range instances {
		instance := strings.ToLower(instanceName)
		current := iscenv.GetInstances()
		existing := current.Find(instance)

		if existing != nil {
			validInstances = append(validInstances, existing)
		} else {
			fmt.Printf("No such instance, name: %s\n", instanceName)
		}
	}

	if len(validInstances) == 0 {
		iscenv.Fatalf("No valid instances provided")
	}

	iscenv.Ensure(iscenv.WithInstance(validInstances[0], copyCSPGateway))
	iscenv.Ensure(configureModCSP)

	for _, instance := range validInstances {
		iscenv.Ensure(iscenv.WithInstance(instance, createApacheSite))
	}

	iscenv.Ensure(func() error { return configureCSPGateway(validInstances) })
	iscenv.Ensure(restartApache)
}

func copyCSPGateway(instance *iscenv.ISCInstance) error {
	var err error

	_, err = os.Stat(gatewayDir)
	if !os.IsNotExist(err) {
		fmt.Printf("CSP Gateway directory already exists, skipping copy, path: %s\n", gatewayDir)
		return err
	}

	fmt.Printf("Copying CSP Gateway from instance: %s, path: %s\n", instance.Name, gatewayDir)

	if err := os.MkdirAll(gatewayDir, 0755); err != nil {
		return err
	}

	if err := iscenv.DockerCopy(instance.Name, srcGatewayDir, gatewayDir); err != nil {
		return err
	}

	// remove the CSP.ini so it can be "replaced" only the first time
	if err := os.Remove(cspIniPath); err != nil {
		return err
	}

	return nil
}

func configureModCSP() error {
	fmt.Println("Configuring Apache CSP module")
	if err := writeTemplate(filepath.Join(apacheDir, "mods-available", "csp.conf"), iscenv.CSPConf, iscenv.ApacheTemplateData{GatewayDir: gatewayDir}, true); err != nil {
		return err
	}

	if err := writeTemplate(filepath.Join(apacheDir, "mods-available", "csp.load"), iscenv.CSPLoad, iscenv.ApacheTemplateData{GatewayDir: gatewayDir}, true); err != nil {
		return err
	}

	out, err := exec.Command("a2enmod", "csp").CombinedOutput()
	fmt.Println(strings.TrimSpace(string(out)))
	if err != nil {
		return err
	}

	return nil
}

func createApacheSite(instance *iscenv.ISCInstance) error {
	siteName := strings.ToLower(instance.Name) + "-iscenv"
	fmt.Printf("Creating Apache site, name: %s\n", siteName)

	if err := writeTemplate(filepath.Join(apacheDir, "sites-available", siteName+".conf"), iscenv.SiteConf, iscenv.ApacheTemplateData{Instance: instance}, false); err != nil {
		return err
	}

	out, err := exec.Command("a2ensite", siteName).CombinedOutput()
	fmt.Println(strings.TrimSpace(string(out)))
	if err != nil {
		return err
	}

	return nil
}

func configureCSPGateway(instances iscenv.ISCInstances) error {
	fmt.Println("Configuring CSP Gateway")
	return writeTemplate(cspIniPath, iscenv.CSPIni, iscenv.ApacheTemplateData{Instances: instances}, false)
}

func restartApache() error {
	fmt.Println("Restarting Apache")
	out, err := exec.Command("service", "apache2", "restart").CombinedOutput()
	fmt.Println(strings.TrimSpace(string(out)))
	if err != nil {
		return err
	}

	return nil
}

func writeTemplate(path string, tmpl *template.Template, data interface{}, overwrite bool) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if err != nil {
			return err
		}

		if !overwrite {
			fmt.Println("WARNING: Template already exists, will write .new file")
			path = path + ".new"
		}
	}

	fmt.Printf("Writing template, path: %s\n", path)

	w, err := os.Create(path)
	if err != nil {
		return err
	}

	defer w.Close()
	return tmpl.Execute(w, data)
}
