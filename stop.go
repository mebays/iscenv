/*
Copyright 2015 Ontario Systems

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

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var stopTimeout uint

var stopCommand = &cobra.Command{
	Use:   "stop [OPTIONS] INSTANCE [INSTANCE...]",
	Short: "Stop an ISC product instance",
	Long:  "Stop a running ISC product instance, attempting a safe shutdown",
}

func init() {
	stopCommand.Run = stop
	stopCommand.Flags().UintVarP(&stopTimeout, "time", "t", 60, "The amount of time to wait for the instance to stop cleanly before killing it.")
	addMultiInstanceFlags(stopCommand, "stop")
}

func stop(_ *cobra.Command, args []string) {
	instances := multiInstanceFlags.getInstances(args)
	for _, instanceName := range instances {
		instance := strings.ToLower(instanceName)
		current := getInstances()
		existing := current.find(instance)

		if existing != nil {
			err := dockerClient.StopContainer(existing.ID, stopTimeout)
			if err != nil {
				fatalf("Could not stop instance, name: %s, error: %s\n", existing.Name, err)
			}

			fmt.Println(existing.ID)
		} else {
			fmt.Printf("No such instance, name: %s\n", instanceName)
		}
	}
}
