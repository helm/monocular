/*
Copyright (c) 2018 The Helm Authors

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

package utils

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version          = "devel"
	UserAgentComment string
)

// UserAgent returns the user agent to be used during calls to the chart repositories
// Examples:
// chart-repo/devel
// chart-repo/1.0
// chart-repo/1.0 (monocular v1.0-beta4)
// More info here https://github.com/kubeapps/kubeapps/issues/767#issuecomment-436835938
func UserAgent() string {
	ua := "chart-repo/" + Version
	if UserAgentComment != "" {
		ua = fmt.Sprintf("%s (%s)", ua, UserAgentComment)
	}
	return ua
}

//VersionCmd returns Monocular version information
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "returns version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
