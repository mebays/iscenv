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

package isclib

// TODO: Consider making a pass through this entire library and using errwrap where appropriate

import (
	"bufio"
	"bytes"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

const (
	defaultCControlPath = "ccontrol"
	defaultCSessionPath = "csession"
	cacheParametersFile = "parameters.isc"
)

const (
	importXMLHeader = `<?xml version="1.0" encoding="UTF-8"?>
<Export generator="Cache" version="25">
<Routine name="%s" type="MAC" languagemode="0"><![CDATA[
EnsLibMain() public {
	try {
		do MAIN
	} catch ex {
		do BACK^%%ETN
		use 0
		write !,"Exception: ",ex.DisplayString(),!,"  name: ",ex.Name,!,"  code: ",ex.Code,!
		do $zutil(4, $job, 99)
	}
}

`
	importXMLFooter = `
]]></Routine>
</Export>`
)

var (
	globalCControlPath        = defaultCControlPath
	globalCSessionPath        = defaultCSessionPath
	executeTemporaryDirectory = "" // Default is system temp directory
)

// CControlPath returns the current path to the ccontrol executable
func CControlPath() string { return globalCControlPath }

// SetCControlPath sets the current path to the ccontrol executable
func SetCControlPath(path string) {
	globalCControlPath = path
}

// CSessionPath returns the current path to the csession executable
func CSessionPath() string { return globalCSessionPath }

// SetCSessionPath sets the current path to the csession executable
func SetCSessionPath(path string) {
	globalCSessionPath = path
}

// ExecuteTemporaryDirectory returns the directory where temporary files for ObjectScript execution will be placed.
// "" means the system default temp directory.
func ExecuteTemporaryDirectory() string {
	return executeTemporaryDirectory
}

// SetExecuteTemporaryDirectory sets the directory where temporary files for ObjectScript execution will be placed.
// Passing "" will result in using the system default temp directory.
func SetExecuteTemporaryDirectory(path string) {
	executeTemporaryDirectory = path
}

// ISCExists returns a boolean which is true when an ISC product or Caché instance exists on this system.
func ISCExists() bool {
	if _, err := exec.LookPath(globalCControlPath); err != nil {
		log.WithField("ccontrolPath", globalCControlPath).WithError(err).Debug("ccontrol executable not found")
		return false
	}

	if _, err := exec.LookPath(globalCSessionPath); err != nil {
		log.WithField("csessionPath", globalCControlPath).WithError(err).Debug("csession executable not found")
		return false
	}

	return true
}

// LoadInstances returns a listing of all Caché/Ensemble instances on this system.
// It returns the list of instances and any error encountered.
func LoadInstances() (Instances, error) {
	qs, err := qlist("")
	if err != nil {
		return nil, err
	}

	instances := make(Instances, 0)
	scanner := bufio.NewScanner(bytes.NewBufferString(qs))
	for scanner.Scan() {
		instance, err := InstanceFromQList(scanner.Text())
		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return instances, nil
}

// LoadInstance retrieves a single instance by name.
// The instance name is case insensitive.
// It returns the instance and any error encountered.
func LoadInstance(name string) (*Instance, error) {
	q, err := qlist(name)
	if err != nil {
		return nil, err
	}
	return InstanceFromQList(q)
}

// InstanceFromQList will parse the output of a ccontrol qlist into an Instance struct.
// It expects the results of a ccontrol qlist for a single instance as a string.
// It returns the parsed instance and any error encountered.
func InstanceFromQList(qlist string) (*Instance, error) {
	i := new(Instance)
	if err := i.UpdateFromQList(qlist); err != nil {
		return nil, err
	}

	return i, nil
}
