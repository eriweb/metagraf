/*
Copyright 2019 The MetaGraph Authors

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
	"github.com/golang/glog"
	log "k8s.io/klog"
	"metagraf/mg/cmd"
	"metagraf/pkg/mgver"
	"os"
)


func main() {
	if len(mgver.GitTag) == 0 {
		cmd.MGVersion = mgver.GitBranch+"("+mgver.GitHash+")"
	} else {
		cmd.MGVersion = mgver.GitTag+"("+mgver.GitHash+")"
	}

	err := cmd.Execute()
	if err != nil {
		glog.Error(err)
		glog.Flush()
		os.Exit(1)
	}
	log.Flush()
	os.Exit(0)
}
