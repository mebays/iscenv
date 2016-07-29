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

import "github.com/ontariosystems/iscenv/internal/app"

// getActivatedLifecyclers will populate lifecyclers with activated lifecyclers plugins based on the provided list.
// It returns the close function from the manager so you can easily defer the return.  It will log fatally on any errors
func getActivatedLifecyclers(pluginsToActivate []string, args app.PluginArgs, lifecyclers *[]*app.ActivatedLifecycler) func() {
	var err error
	lcm, err := app.NewLifecyclerManager(args)
	if err != nil {
		app.ErrorLogger(nil, err).Fatal("Failed to create lifecycle plugin manager")
	}

	*lifecyclers, err = lcm.ActivatePlugins(pluginsToActivate)
	if err != nil {
		app.ErrorLogger(nil, err).Fatal("Failed to activate lifecycle plugins")
	}

	return lcm.Close
}