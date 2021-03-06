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

import (
	"fmt"
	"strings"
)

// NewPluginFlags creates and returns an empty map of PluginFlag
func NewPluginFlags() PluginFlags {
	return PluginFlags{
		Flags: make(map[string]*PluginFlag),
	}
}

// PluginFlags is a slice of PluginFlag
type PluginFlags struct {
	Flags map[string]*PluginFlag
}

// AddFlag adds a Plugin Flag to the list of available flags.
func (pf *PluginFlags) AddFlag(flag string, hasConfig bool, defaultValue interface{}, usage string) error {
	flag = strings.ToLower(flag)
	if _, ok := pf.Flags[flag]; ok {
		return fmt.Errorf("Flag already exists, flag: %s", flag)
	}

	pf.Flags[flag] = NewPluginFlag(flag, hasConfig, defaultValue, usage)
	return nil
}
