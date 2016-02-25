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

package app

import (
	"strconv"

	"github.com/fsouza/go-dockerclient"
	"github.com/ontariosystems/iscenv/iscenv"
)

func DockerPort(port iscenv.ContainerPort) docker.Port {
	return docker.Port(port.String()) + "/tcp"
}

func DockerPortBinding(port int64, portOffset int64) []docker.PortBinding {
	return []docker.PortBinding{docker.PortBinding{HostIP: "", HostPort: strconv.FormatInt(port+portOffset, 10)}}
}

// Assumes a single binding
func GetDockerBindingPort(bindings []docker.PortBinding) (iscenv.ContainerPort, error) {
	port, err := strconv.ParseInt(bindings[0].HostPort, 10, 64)
	if err != nil {
		return 0, err
	}

	return iscenv.ContainerPort(port), nil
}
