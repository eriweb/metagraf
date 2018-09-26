/*
Copyright 2018 The MetaGraph Authors

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
	"metagraf/pkg/metagraf"
	"os"
	"strings"
)

// Returns a slice of strings of potential parameterized variables in a
// metaGraf specification that can be found in the execution environment.
func VarsFromEnv(mgv MGVars) EnvVars {
	envs := EnvVars{}
	for _,v := range os.Environ() {
		key, val := keyValueFromEnv(v)
		if _, ok := mgv[key]; ok {
			envs[key] = val
		}
	}
	return envs
}

func VarsFromCmd(mgv MGVars, cvars CmdVars) map[string]string {
	vars := make(map[string]string)

	for k,v := range cvars {
		if _, ok := mgv[k]; ok {
			vars[k] = v
		}
	}
	return vars
}

func keyValueFromEnv(s string) (string,string) {
	return strings.Split(s,"=")[0],strings.Split(s,"=")[1]
}

// Returns a slice of strings of alle parameterized fields in a metaGraf
// specification.
// @todo need to look for parameterized fields in more places
func VarsFromMetaGraf(mg *metagraf.MetaGraf) MGVars {
	vars := MGVars{}

	// Environment Section
	for _,env := range mg.Spec.Environment.Local {
		vars[env.Name] = ""
	}
	for _,env := range mg.Spec.Environment.External.Introduces {
		vars[env.Name] = ""
	}
	for _,env := range mg.Spec.Environment.External.Consumes {
		vars[env.Name] = ""
	}

	// Config section

	return vars
}

// Returns a list of variables from command line or environment where
// command line is the most significant.
func OverrideVars(mg *metagraf.MetaGraf, cvars CmdVars) map[string]string {
	ovars := make(map[string]string)

	// Fetch possible variables form metaGraf specification
	mgvars := VarsFromMetaGraf(mg)
	for k,v := range VarsFromEnv(mgvars) {
		ovars[k] = v
	}
	for k,v := range VarsFromCmd(mgvars,cvars) {
		ovars[k] = v
	}

	return ovars
}