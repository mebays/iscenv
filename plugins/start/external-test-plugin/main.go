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

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ontariosystems/iscenv/iscenv"
)

var plog = log.WithField("plugin", "external-test-plugin")

type Plugin struct{}

func main() {
	iscenv.ServeStartPlugin(new(Plugin))
}

func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	plog.Info("Flags")
	return iscenv.NewPluginFlags(), nil
}

func (*Plugin) Environment(version string, flags map[string]interface{}) ([]string, error) {
	plog.Info("Environment")
	return nil, nil
}

func (*Plugin) Volumes(_ string, _ map[string]interface{}) ([]string, error) {
	plog.Info("Volumes")
	return nil, nil
}

func (*Plugin) Ports(_ string, _ map[string]interface{}) ([]string, error) {
	plog.Info("Ports")
	return nil, nil
}

func (*Plugin) BeforeInstance(state iscenv.InternalInstance) error {
	plog.Info("BeforeInstance")
	return nil
}

func (*Plugin) WithInstance(state iscenv.InternalInstance) error {
	plog.Info("WithInstance")
	return nil
}

func (*Plugin) AfterInstance(state iscenv.InternalInstance) error {
	plog.Info("AfterInstance")
	return nil
}
