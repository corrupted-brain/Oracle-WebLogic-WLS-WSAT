// Copyright © 2018 Kevin Kirsche <kev.kirsche[at]gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// BuildVersion the version of the binary (e.g. 1.2.0)
	BuildVersion string
	// BuildGoVersion is what version of Golang we built the binary with
	BuildGoVersion string
	// BuildHash is the git hash that we were at when this was built
	BuildHash string
	// BuildTime is when we built the binary
	BuildTime string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "The version of the binary",
	Long: `The build date and build hash associated with the build to allow for
	better identification of when the binary was made and what features it
	offers`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Printf("Version:\t%s", BuildVersion)
		logrus.Printf("Go Version:\t%s", strings.Join(strings.Split(BuildGoVersion, "|^|"), " "))
		logrus.Printf("Git Hash:\t%s", BuildHash)
		logrus.Printf("Build Time:\t%s", BuildTime)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
