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

// Instances represents a collection of Caché/Ensemble instances
type Instances []*Instance

// Update will query the underlying instances and update the Instance fields with their current values.
// It returns any error encountered.
func (instances Instances) Update() error {
	for _, instance := range instances {
		if err := instance.Update(); err != nil {
			return err
		}
	}

	return nil
}
