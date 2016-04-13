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

package iscenv

type ISCVersions []*ISCVersion

func (evs ISCVersions) Latest() *ISCVersion {
	return evs[len(evs)-1]
}

func (evs ISCVersions) Exists(versionString string) bool {
	return evs.Find(versionString) != nil
}

func (evs ISCVersions) Find(versionString string) *ISCVersion {
	for _, version := range evs {
		if version.Version == versionString {
			return version
		}
	}

	return nil
}